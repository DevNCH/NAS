package auth

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// Session representa uma sessão autenticada de um usuário.
type Session struct {
	UserID    int
	Username  string
	Role      string
	ExpiresAt time.Time
}

// SessionManager guarda as sessões ativas em memória.
// Simples e suficiente para um servidor NAS de rede local;
// se o processo reiniciar, todas as sessões expiram (usuários
// precisam logar novamente).
type SessionManager struct {
	mu       sync.RWMutex
	sessions map[string]Session
	ttl      time.Duration
}

// NewSessionManager cria um gerenciador com um tempo de vida (ttl)
// padrão para cada sessão criada.
func NewSessionManager(ttl time.Duration) *SessionManager {
	sm := &SessionManager{
		sessions: make(map[string]Session),
		ttl:      ttl,
	}

	go sm.cleanupLoop()

	return sm
}

// Create gera um novo token de sessão para o usuário informado.
func (sm *SessionManager) Create(userID int, username string, role string) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	sm.mu.Lock()
	sm.sessions[token] = Session{
		UserID:    userID,
		Username:  username,
		Role:      role,
		ExpiresAt: time.Now().Add(sm.ttl),
	}
	sm.mu.Unlock()

	return token, nil
}

// Get retorna a sessão associada ao token, se ela existir e ainda for válida.
func (sm *SessionManager) Get(token string) (Session, bool) {
	sm.mu.RLock()
	session, ok := sm.sessions[token]
	sm.mu.RUnlock()

	if !ok {
		return Session{}, false
	}

	if time.Now().After(session.ExpiresAt) {
		sm.Delete(token)
		return Session{}, false
	}

	return session, true
}

// Delete invalida um token de sessão (usado no logout).
func (sm *SessionManager) Delete(token string) {
	sm.mu.Lock()
	delete(sm.sessions, token)
	sm.mu.Unlock()
}

// cleanupLoop remove periodicamente sessões expiradas da memória.
func (sm *SessionManager) cleanupLoop() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()

		sm.mu.Lock()
		for token, session := range sm.sessions {
			if now.After(session.ExpiresAt) {
				delete(sm.sessions, token)
			}
		}
		sm.mu.Unlock()
	}
}

// generateToken cria um token aleatório seguro para identificar a sessão.
func generateToken() (string, error) {
	buf := make([]byte, 32)

	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}

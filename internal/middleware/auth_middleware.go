package middleware

import (
	"net/http"

	"github.com/DevNCH/NAS/internal/auth"

	"github.com/gin-gonic/gin"
)

const (
	// SessionCookieName é o nome do cookie que guarda o token de sessão.
	SessionCookieName = "session_token"

	// ContextUserKey é a chave usada para guardar a sessão no contexto do Gin.
	ContextUserKey = "user"
)

// RequireAuth garante que a requisição tenha um cookie de sessão válido.
// Se não tiver, interrompe a requisição com 401.
func RequireAuth(sessions *auth.SessionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(SessionCookieName)
		if err != nil || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "não autenticado",
			})
			return
		}

		session, ok := sessions.Get(token)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "sessão inválida ou expirada",
			})
			return
		}

		// Disponibiliza a sessão do usuário para os handlers seguintes.
		c.Set(ContextUserKey, session)

		c.Next()
	}
}

// RequireRole garante que o usuário autenticado tenha um dos papéis
// (roles) permitidos. Deve ser usado depois de RequireAuth.
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}

	return func(c *gin.Context) {
		value, exists := c.Get(ContextUserKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "não autenticado",
			})
			return
		}

		session, ok := value.(auth.Session)
		if !ok || !allowed[session.Role] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "acesso negado para este perfil",
			})
			return
		}

		c.Next()
	}
}

// CurrentUser é um helper para os handlers pegarem a sessão do usuário
// autenticado a partir do contexto do Gin.
func CurrentUser(c *gin.Context) (auth.Session, bool) {
	value, exists := c.Get(ContextUserKey)
	if !exists {
		return auth.Session{}, false
	}

	session, ok := value.(auth.Session)
	return session, ok
}

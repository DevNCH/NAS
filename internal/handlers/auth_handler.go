package handlers

import (
	"net/http"

	"github.com/DevNCH/NAS/internal/auth"
	"github.com/DevNCH/NAS/internal/middleware"
	"github.com/DevNCH/NAS/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
	sessions    *auth.SessionManager
}

func NewAuthHandler(authService *services.AuthService, sessions *auth.SessionManager) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		sessions:    sessions,
	}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login valida usuário/senha e, se corretos, cria uma sessão e devolve
// o token em um cookie HTTP-only.
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "usuário e senha são obrigatórios"})
		return
	}

	user, err := h.authService.Login(req.Username, req.Password)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário ou senha inválidos"})
		return
	}

	token, err := h.sessions.Create(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao criar sessão"})
		return
	}

	// secure=false porque o servidor roda em rede local (TV Box) sem HTTPS.
	// Se o projeto passar a rodar atrás de HTTPS, mude para true.
	c.SetCookie(
		middleware.SessionCookieName,
		token,
		int((24 * 3600)), // 24 horas, em segundos
		"/",
		"",
		false,
		true, // HttpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}

// Logout invalida a sessão atual e remove o cookie.
func (h *AuthHandler) Logout(c *gin.Context) {
	if token, err := c.Cookie(middleware.SessionCookieName); err == nil {
		h.sessions.Delete(token)
	}

	c.SetCookie(middleware.SessionCookieName, "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "logout realizado"})
}

// Me devolve os dados do usuário autenticado. Rota protegida, útil para
// o frontend verificar se a sessão ainda é válida.
func (h *AuthHandler) Me(c *gin.Context) {
	session, ok := middleware.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "não autenticado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       session.UserID,
		"username": session.Username,
		"role":     session.Role,
	})
}

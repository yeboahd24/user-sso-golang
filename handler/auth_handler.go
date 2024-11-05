package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yeboahd24/user-sso/config"
	"github.com/yeboahd24/user-sso/service"
	"github.com/yeboahd24/user-sso/util"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthHandler struct {
	authService *service.AuthService
	config      *config.Config
}

func NewAuthHandler(authService *service.AuthService, config *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		config:      config,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input struct {
		Email    string `form:"email" binding:"required,email"`
		Password string `form:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	if err := h.authService.RegisterUser(input.Email, input.Password); err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	// Redirect to login page or home page after successful registration
	c.Redirect(http.StatusSeeOther, "/")
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	url := h.authService.GetGoogleAuthURL()
	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
		return
	}

	user, err := h.authService.HandleGoogleCallback(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT token here
	token, err := util.GenerateJWT(user, h.config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	// Log the content type and raw data
	fmt.Printf("Content-Type: %s\n", c.GetHeader("Content-Type"))
	body, _ := io.ReadAll(c.Request.Body)
	fmt.Printf("Raw request body: %s\n", string(body))
	// Restore the body for binding
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// Try binding JSON first
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("JSON binding error: %v\n", err)
		// If JSON binding fails, try form data
		if err := c.ShouldBind(&req); err != nil {
			fmt.Printf("Form binding error: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error":         fmt.Sprintf("Binding error: %v", err),
				"received_data": string(body),
			})
			return
		}
	}

	// Log the parsed request
	fmt.Printf("Parsed request: %+v\n", req)

	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	user, err := h.authService.ValidateLogin(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate tokens
	token, err := h.authService.GenerateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set cookie
	c.SetCookie(
		"session_token",
		token,
		int(time.Hour*24), // 24 hours
		"/",
		"",
		false, // set to true in production
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

func (h *AuthHandler) GetGoogleAuthURL(c *gin.Context) {
	url := h.authService.GetGoogleAuthURL()
	c.JSON(http.StatusOK, gin.H{"url": url})
}

// VerifySession checks if the current session is valid
func (h *AuthHandler) VerifySession(c *gin.Context) {
	token, err := c.Cookie("session_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No session found"})
		return
	}

	// Validate token
	claims, err := h.authService.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Session is valid",
		"user_id": claims.UserID,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Clear the session cookie
	c.SetCookie(
		"session_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

package handler

import (
	"auth-service/internal/dto"
	"auth-service/internal/service"
	"auth-service/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func NewAuthHandler(
	authService *service.AuthService,
) *AuthHandler {

	return &AuthHandler{
		AuthService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {

	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	err := utils.Validate.Struct(req)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": utils.FormatValidationError(err),
		})

		return
	}

	err = h.AuthService.Register(req)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "register success",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {

	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	err := utils.Validate.Struct(req)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": utils.FormatValidationError(err),
		})

		return
	}

	res, err := h.AuthService.Login(req)

	if err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) Refresh(c *gin.Context) {

	var req dto.RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	res, err := h.AuthService.RefreshToken(req)

	if err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) Logout(c *gin.Context) {

	var req dto.RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	refresh, err := h.AuthService.RefreshRepo.
		FindByToken(req.RefreshToken)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid refresh token",
		})

		return
	}

	err = h.AuthService.RefreshRepo.
		Revoke(refresh.ID.String())

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "logout failed",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "logout success",
	})
}

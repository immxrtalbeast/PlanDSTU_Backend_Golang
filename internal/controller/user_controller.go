package controller

import (
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/immxrtalbeast/plandstu/internal/domain"
)

type UserController struct {
	interactor  domain.UserInteractor
	tokenTTL    time.Duration
	tokenSecret string
}

func NewUserController(interactor domain.UserInteractor, tokenTTL time.Duration, tokenSecret string) *UserController {
	return &UserController{interactor: interactor, tokenTTL: tokenTTL, tokenSecret: tokenSecret}
}

func (c *UserController) Register(ctx *gin.Context) {
	type RegisterRequest struct {
		Login string `json:"login" binding:"required,min=3,max=50"`
		Pass  string `json:"password" binding:"required,min=8,max=50"`
		Group string `json:"group"`
	}

	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Валидация пароля
	passRegex := regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+\[\]{};:<>,./?~\\-]+$`)
	if !passRegex.MatchString(req.Pass) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid password",
			"details": "Password contains forbidden characters",
		})
		return
	}

	// Если все проверки пройдены
	id, err := c.interactor.CreateUser(ctx, req.Login, req.Pass, req.Group)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to create user",
			"details": err.Error(),
		})
		return
	}
	token, err := c.interactor.Login(ctx, req.Login, req.Pass)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to login",
			"details": err.Error(),
		})
		return
	}
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		"jwt",                     // Имя куки
		token,                     // Значение токена
		int(c.tokenTTL.Seconds()), // Макс возраст в секундах
		"/",                       // Путь
		"",                        // Домен (пусто для текущего домена)
		false,                     // Secure (использовать true в production для HTTPS)
		false,                     // HttpOnly
	)

	// Устанавливаем куку
	ctx.SetCookie(
		"user_id",
		id.String(),
		int(c.tokenTTL.Seconds()),
		"/",
		"",
		false,
		false, // HttpOnly=false чтобы клиент мог читать JS
	)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user created",
	})

}
func (c *UserController) Login(ctx *gin.Context) {
	type LoginRequest struct {
		Login string `json:"login" binding:"required"`
		Pass  string `json:"password" binding:"required"`
	}
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}
	token, err := c.interactor.Login(ctx, req.Login, req.Pass)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to login",
			"details": err.Error(),
		})
		return
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		"jwt",                     // Имя куки
		token,                     // Значение токена
		int(c.tokenTTL.Seconds()), // Макс возраст в секундах
		"/",                       // Путь
		"",                        // Домен (пусто для текущего домена)
		false,                     // Secure (использовать true в production для HTTPS)
		false,                     // HttpOnly
	)

	ctx.JSON(http.StatusOK, gin.H{})
}

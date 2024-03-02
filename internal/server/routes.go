package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ostheperson/go-auth-service/integrations"
	"github.com/ostheperson/go-auth-service/internal/domain"
	"github.com/ostheperson/go-auth-service/internal/handlers"
	"github.com/ostheperson/go-auth-service/internal/services"
)

func RegisterRoutes(s *domain.Server) http.Handler {
	r := gin.Default()

	hh := NewHelloHandler(s)
	r.GET("/", hh.HelloWorldHandler)
	r.GET("/health", hh.healthHandler)
	const (
		authRoute = "/auth"
		userRoute = "/users"
	)
	mailer := integrations.NewMailerClient(s.Env.POSTMARK_API_KEY)
	tokenService := services.NewTokenService(s.Db.GetClient())
	ah := handlers.NewAuthHandler(s, mailer, tokenService)
	uh := handlers.NewUsersHandler(s)

	authRoutes := r.Group(authRoute)
	{
		authRoutes.POST("/user/signup", ah.SignUp)
		authRoutes.POST("/user/signin", ah.SignIn)
		authRoutes.POST("/admin/signin", ah.SignIn)
		authRoutes.POST("/refresh", ah.RefreshToken)
		authRoutes.POST("/email-verify/request", ah.ResendVerifyEmail)
		authRoutes.POST("/email-verify/confirm", ah.ConfirmEmail)
		authRoutes.POST("/reset-password/request", ah.ForgotPassword)
		authRoutes.POST("/reset-password/confirm", ah.ResetPassword)
	}

	// USERS
	userRoutes := r.Group(userRoute)
	userRoutes.Use(JwtAuthMiddleware(s.Env.AccessTokenSecret))
	{
		userRoutes.GET("", RoleMiddleware(domain.AdminRole), uh.GetUsers)
		userRoutes.GET(
			"/:id",
			RoleMiddleware(domain.AdminRole, domain.UserRole),
			uh.GetUser,
		)
		userRoutes.PATCH(
			"/:id",
			RoleMiddleware(domain.AdminRole, domain.UserRole),
			uh.UpdateUser,
		)
		userRoutes.DELETE(
			"/:id",
			RoleMiddleware(domain.AdminRole),
			uh.RemoveUser,
		)
		userRoutes.DELETE("/all", RoleMiddleware(domain.AdminRole))
	}

	return r
}

type HelloHandler struct {
	*domain.Server
}

func NewHelloHandler(s *domain.Server) *HelloHandler {
	return &HelloHandler{Server: s}
}

func (s *HelloHandler) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "live"

	c.JSON(http.StatusOK, resp)
}

func (s *HelloHandler) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.Db.Health())
}

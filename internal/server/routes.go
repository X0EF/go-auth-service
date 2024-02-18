package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ostheperson/go-auth-service/internal/domain"
	"github.com/ostheperson/go-auth-service/internal/handlers"
)

func RegisterRoutes(s *domain.Server) http.Handler {
	r := gin.Default()

	hh := NewHelloHandler(s)
	r.GET("/", hh.HelloWorldHandler)
	r.GET("/health", hh.healthHandler)

	ah := handlers.NewAuthHandler(s)
	r.POST("/auth/user/signup", ah.SignUp)
	r.POST("/auth/user/signin", ah.SignIn)
	r.POST("/auth/admin/signin", ah.SignIn)

	protected := r.Group("")
	protected.Use(JwtAuthMiddleware(s.Env.AccessTokenSecret))
	uh := handlers.NewUsersHandler(s)

	protected.GET("/users", RoleMiddleware(domain.AdminRole), uh.GetUsers)
	protected.GET("/users/:id", RoleMiddleware(domain.AdminRole, domain.UserRole), uh.GetUser)
	protected.PATCH("/users/:id", RoleMiddleware(domain.AdminRole, domain.UserRole), uh.UpdateUser)
	protected.DELETE("/users/:id", RoleMiddleware(domain.AdminRole, domain.UserRole), uh.RemoveUser)
	protected.DELETE("/users/all", RoleMiddleware(domain.AdminRole))

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

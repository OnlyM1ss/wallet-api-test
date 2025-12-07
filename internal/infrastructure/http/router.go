package http

import (
	"github.com/gin-gonic/gin"
	"walletapitest/internal/infrastructure/http/handlers"
	"walletapitest/internal/infrastructure/http/middlewares"
)

type Router struct {
	engine         *gin.Engine
	userHandler    *handlers.UserHandler
	authMiddleware gin.HandlerFunc
}

func NewRouter(userHandler *handlers.UserHandler) *Router {
	router := &Router{
		engine:      gin.Default(),
		userHandler: userHandler,
	}
	
	router.setupRoutes()
	return router
}

func (r *Router) setupRoutes() {
	// Public routes
	public := r.engine.Group("/api/v1")
	{
		public.POST("/users", r.userHandler.CreateUser)
		public.POST("/login", r.userHandler.Login)
	}
	
	// Protected routes
	protected := r.engine.Group("/api/v1")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/users/:id", r.userHandler.GetUser)
		// другие защищенные маршруты
	}
	
	// Health check
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
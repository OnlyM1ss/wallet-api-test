package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"walletapitest/internal/config"
	"walletapitest/internal/domain/services"
	postgres "walletapitest/internal/infrastructure/database/postgres"
	"walletapitest/internal/infrastructure/http/handlers"
	"walletapitest/internal/pkg/logger"
)

type App struct {
	cfg    *config.Config
	logger logger.Logger
	router *gin.Engine
	db     *sqlx.DB
}

func New(cfg *config.Config, logger logger.Logger) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *App) Run() error {
	// Инициализация базы данных
	db, err := a.initDB()
	if err != nil {
		a.logger.Error("Failed to connect to database", "error", err)
		// В режиме разработки можно продолжить без БД
		a.db = nil
	} else {
		a.db = db
		defer db.Close()
	}

	// Инициализация зависимостей
	var userService *services.UserService
	var walletService *services.WalletService
	if a.db != nil {
		// Инициализация репозиториев
		userRepo := postgres.NewUserRepository(db)
		walletRepo := postgres.NewWalletRepository(db)

		// Инициализация сервисов
		userService = services.NewUserService(userRepo)
		walletService = services.NewWalletService(walletRepo)
	} else {
		// База данных обязательна для работы приложения
		return err
	}

	// Инициализация хендлеров
	userHandler := handlers.NewUserHandler(userService)
	walletHandler := handlers.NewWalletHandler(walletService)

	// Инициализация роутера
	a.router = a.initRouter(userHandler, walletHandler)

	// Запуск сервера
	srv := &http.Server{
		Addr:         ":" + a.cfg.Server.Port,
		Handler:      a.router,
		ReadTimeout:  time.Duration(a.cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(a.cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(a.cfg.Server.IdleTimeout) * time.Second,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("Failed to start server", "error", err)
		}
	}()

	a.logger.Info("Server started on port " + a.cfg.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		a.logger.Error("Server forced to shutdown", "error", err)
		return err
	}

	a.logger.Info("Server exited properly")
	return nil
}

func (a *App) initRouter(userHandler *handlers.UserHandler, walletHandler *handlers.WalletHandler) *gin.Engine {
	router := gin.Default()

	// Public routes
	router.POST("/api/v1/users", userHandler.CreateUser)
	router.POST("/api/v1/login", userHandler.Login)
	router.GET("/api/v1/users/:id", userHandler.GetUser)

	// Wallet routes
	router.POST("/api/v1/wallet", walletHandler.ProcessOperation)
	router.GET("/api/v1/wallet/:walletId", walletHandler.GetWallet)
	router.POST("/api/v1/wallet/create", walletHandler.CreateWallet)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		if a.db != nil {
			// Проверяем соединение с БД
			if err := a.db.Ping(); err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": err.Error()})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	return router
}

func (a *App) initDB() (*sqlx.DB, error) {
	connStr := "postgres://" + a.cfg.Database.User + ":" + a.cfg.Database.Password +
		"@" + a.cfg.Database.Host + ":" + a.cfg.Database.Port + "/" + a.cfg.Database.Name +
		"?sslmode=" + a.cfg.Database.SSLMode

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Проверка соединения
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Создание таблицы users если её нет
	if err := a.createUsersTable(db); err != nil {
		a.logger.Error("Failed to create users table", "error", err)
		return nil, err
	}

	// Создание таблицы wallets если её нет
	if err := a.createWalletsTable(db); err != nil {
		a.logger.Error("Failed to create wallets table", "error", err)
		return nil, err
	}

	a.logger.Info("Database connection established")
	return db, nil
}

func (a *App) createUsersTable(db *sqlx.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			username VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`

	_, err := db.Exec(query)
	return err
}

func (a *App) createWalletsTable(db *sqlx.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS wallets (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			balance BIGINT NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX IF NOT EXISTS idx_wallets_user_id ON wallets(user_id);
		CREATE INDEX IF NOT EXISTS idx_wallets_id ON wallets(id);
	`

	_, err := db.Exec(query)
	return err
}

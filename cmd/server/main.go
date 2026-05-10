package main

import (
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"

	pb "github.com/novaherdi/go-proto/gen/guestbook/v1"
	"github.com/novaherdi/go-proto/internal/controller"
	"github.com/novaherdi/go-proto/internal/di"
	"github.com/novaherdi/go-proto/internal/model"
	"github.com/novaherdi/go-proto/internal/repository"
	"github.com/novaherdi/go-proto/internal/service"
	"github.com/novaherdi/go-proto/pkg/auth"
	"github.com/novaherdi/go-proto/pkg/config"
	"github.com/novaherdi/go-proto/pkg/framework"
	"github.com/novaherdi/go-proto/pkg/router"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 0. Load Configuration using Viper
	// Viper supports:
	// - JSON config files
	// - Environment variables (APP_ prefix, e.g., APP_SERVER_PORT=8000)
	// - Nested keys with underscore (e.g., APP_DATABASE_HOST)
	//
	// Environment variables override file settings.
	// Example: APP_DATABASE_HOST=prod.db.com APP_JWT_SECRETKEY=prod-key ./server
	env := os.Getenv("APP_ENV")
	configFile := "pkg/config/config.local.json"
	if env == "production" {
		configFile = "pkg/config/config.prod.json"
	}

	log.Printf("Loading config from: %s", configFile)
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 1. Initialize the database with GORM
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the schema (including User model)
	if err := db.AutoMigrate(&model.GuestbookEntry{}, &model.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed successfully")

	// Get the underlying SQL DB for connection pooling and health checks
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	defer sqlDB.Close()

	// Verify connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize JWT service
	jwtConfig := auth.JWTConfig{
		SecretKey:     cfg.JWT.SecretKey,
		TokenDuration: cfg.JWT.TokenDuration,
		Issuer:        cfg.JWT.Issuer,
	}
	jwtService := auth.NewJWTService(jwtConfig)

	// Define protected methods that require JWT authentication
	protectedMethods := map[string]bool{
		"/guestbook.v1.GuestbookService/AddEntry": true,
	}

	// Composition root: wire guestbook stack (repository → service → gRPC adapter).
	controllerFactory := di.NewGuestbookControllerFactoryWithCache(db, cfg.Redis)
	guestbookController := controllerFactory.CreateController()

	// Auth stack: user persistence + token issuer abstraction (Dependency Inversion).
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, service.NewJWTTokenIssuer(jwtService))
	authController := controller.NewAuthController(authService)

	// gRPC: logging (outer) → JWT → recovery (inner, adjacent to handler)
	srv := framework.NewServer(
		cfg.Server.Port,
		grpc.ChainUnaryInterceptor(
			framework.UnaryLoggingInterceptor(),
			auth.JWTAuthInterceptor(jwtService, protectedMethods),
			framework.UnaryRecoveryInterceptor(),
		),
	)

	// 3. Register the service with the framework
	srv.RegisterService(func(s *grpc.Server) {
		pb.RegisterGuestbookServiceServer(s, guestbookController)
	})

	// 4. Initialize the HTTP router with auth controller
	httpRouter := router.NewRouter(cfg)
	httpRouter.SetAuthController(authController)
	if err := httpRouter.RegisterGateway(pb.RegisterGuestbookServiceHandlerFromEndpoint); err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}
	httpRouter.SetupRoutes()

	// 5. Run the HTTP Gateway (in a separate goroutine)
	go func() {
		if err := httpRouter.Run(cfg.Server.HttpPort); err != nil {
			log.Fatalf("Failed to serve HTTP gateway: %v", err)
		}
	}()

	// 6. Run the gRPC server
	if err := srv.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

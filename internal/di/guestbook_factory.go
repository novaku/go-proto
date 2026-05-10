// Package di (dependency injection) wires concrete implementations into controllers.
// This is the composition root: application layers stay free of construction logic (SRP).
package di

import (
	"context"
	"fmt"
	"log"

	"github.com/novaherdi/go-proto/internal/controller"
	"github.com/novaherdi/go-proto/internal/repository"
	"github.com/novaherdi/go-proto/internal/service"
	"github.com/novaherdi/go-proto/pkg/cache"
	"github.com/novaherdi/go-proto/pkg/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// GuestbookControllerFactory builds a fully wired GuestbookController for gRPC.
// Open/Closed: swap repository, validator, or cache via factory methods without changing controllers.
type GuestbookControllerFactory struct {
	db    *gorm.DB
	cache cache.Cache
}

// NewGuestbookControllerFactory creates a factory without Redis (cache disabled).
func NewGuestbookControllerFactory(db *gorm.DB) *GuestbookControllerFactory {
	return &GuestbookControllerFactory{
		db:    db,
		cache: nil,
	}
}

// NewGuestbookControllerFactoryWithCache creates a factory and optionally connects Redis.
func NewGuestbookControllerFactoryWithCache(db *gorm.DB, redisConfig config.RedisConfig) *GuestbookControllerFactory {
	var cacheImpl cache.Cache

	if redisConfig.Enabled {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
			Password: redisConfig.Password,
			DB:       redisConfig.DB,
		})

		if err := redisClient.Ping(context.Background()).Err(); err != nil {
			log.Printf("Warning: Failed to connect to Redis: %v", err)
			cacheImpl = nil
		} else {
			log.Println("Redis cache connected successfully")
			cacheImpl = cache.NewRedisCache(redisClient)
		}
	}

	return &GuestbookControllerFactory{
		db:    db,
		cache: cacheImpl,
	}
}

// CreateController builds repository → service → controller with default validator.
func (f *GuestbookControllerFactory) CreateController() *controller.GuestbookController {
	repo := repository.NewGormGuestbookRepository(f.db)
	validator := service.NewDefaultValidator()

	var svc service.GuestbookService
	if f.cache != nil {
		svc = service.NewGuestbookServiceWithCache(repo, validator, f.cache)
	} else {
		svc = service.NewGuestbookService(repo, validator)
	}

	return controller.NewGuestbookController(svc)
}

// CreateControllerWithCustomValidator allows plugging a different GuestbookRequestValidator (Open/Closed).
func (f *GuestbookControllerFactory) CreateControllerWithCustomValidator(validator service.GuestbookRequestValidator) *controller.GuestbookController {
	repo := repository.NewGormGuestbookRepository(f.db)

	var svc service.GuestbookService
	if f.cache != nil {
		svc = service.NewGuestbookServiceWithCache(repo, validator, f.cache)
	} else {
		svc = service.NewGuestbookService(repo, validator)
	}

	return controller.NewGuestbookController(svc)
}

// CreateControllerWithCustomRepository injects a custom GuestbookRepository (testing, alternate stores).
func (f *GuestbookControllerFactory) CreateControllerWithCustomRepository(repo repository.GuestbookRepository) *controller.GuestbookController {
	validator := service.NewDefaultValidator()

	var svc service.GuestbookService
	if f.cache != nil {
		svc = service.NewGuestbookServiceWithCache(repo, validator, f.cache)
	} else {
		svc = service.NewGuestbookService(repo, validator)
	}

	return controller.NewGuestbookController(svc)
}

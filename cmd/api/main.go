package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielbukowski/recipe-app-backend/internal/config"
	"github.com/danielbukowski/recipe-app-backend/internal/recipe"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	_ "github.com/danielbukowski/recipe-app-backend/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title	Recipe API
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfigFromEnvFile()
	if err != nil {
		panic(errors.Join(errors.New("failed to load environment variables"), err))
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(errors.Join(errors.New("failed to create logger"), err))
	}

	pgxCfg, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		panic(errors.Join(errors.New("unable to parse pgx config"), err))
	}

	pgxCfg.MaxConns = 15
	pgxCfg.MaxConnIdleTime = 10
	pgxCfg.MaxConnIdleTime = 15 * time.Minute

	poolCtx, cancelPool := context.WithTimeout(ctx, 5*time.Second)
	defer cancelPool()

	dbpool, err := pgxpool.NewWithConfig(poolCtx, pgxCfg)
	if err != nil {
		panic(errors.Join(errors.New("unable to create connection pool"), err))
	}

	pingCtx, cancelPing := context.WithTimeout(ctx, 5*time.Second)
	defer cancelPing()

	err = dbpool.Ping(pingCtx)
	if err != nil {
		panic(errors.Join(errors.New("failed to ping database"), err))
	}

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	recipeService := recipe.NewService(logger, dbpool)
	recipeHandler := recipe.NewHandler(logger, recipeService)

	recipeHandler.RegisterRoutes(r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HTTPServerPort),
		Handler: r.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("failed to start the http server", zap.Error(err))
			stop()
		}
	}()

	<-ctx.Done()
	fmt.Println("gracefully exiting...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbpool.Close()
	_ = srv.Shutdown(shutdownCtx)

	fmt.Println("closed the application!")
}

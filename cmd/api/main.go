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
	"go.uber.org/zap/zapcore"

	_ "github.com/danielbukowski/recipe-app-backend/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Recipe API
//	@version		0.1
//	@description	A sample of API to recipe backend.
//
//	@host			localhost:8080
//	@BasePath		/api/v1
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	cfg, err := config.LoadConfigFromEnvFile()
	if err != nil {
		panic(errors.Join(errors.New("failed to load environment variables"), err))
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(errors.Join(errors.New("failed to create logger"), err))
	}

	logger.Info("starting the application...")

	pgxCfg, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		panic(errors.Join(errors.New("unable to parse pgx config"), err))
	}

	pgxCfg.MaxConns = 25
	pgxCfg.MinConns = 3
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

	errorLog, err := zap.NewStdLogAt(logger, zapcore.ErrorLevel)
	if err != nil {
		panic(errors.Join(errors.New("failed to create a logger to http errors"), err))
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTPServerPort),
		Handler:      r.Handler(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     errorLog,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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

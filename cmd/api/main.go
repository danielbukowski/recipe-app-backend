package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/danielbukowski/recipe-app-backend/internal/auth"
	"github.com/danielbukowski/recipe-app-backend/internal/cache"
	"github.com/danielbukowski/recipe-app-backend/internal/config"
	"github.com/danielbukowski/recipe-app-backend/internal/healthcheck"
	passwordHasher "github.com/danielbukowski/recipe-app-backend/internal/password-hasher"
	"github.com/danielbukowski/recipe-app-backend/internal/recipe"
	"github.com/danielbukowski/recipe-app-backend/internal/session"

	"github.com/danielbukowski/recipe-app-backend/internal/user"
	"github.com/danielbukowski/recipe-app-backend/internal/validator"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	_ "github.com/danielbukowski/recipe-app-backend/docs"
)

//	@title			Recipe API
//	@version		0.2.0
//	@description	A sample of API to recipe backend.
//
//	@host			localhost:8080
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.LoadEnvironmentVariablesToConfig()
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

	mcache := memcache.New(cfg.MemcachedServer)
	mcache.Timeout = 150 * time.Millisecond

	err = mcache.Ping()
	if err != nil {
		panic(errors.Join(errors.New("failed to ping memcache"), err))
	}

	e := echo.New()
	e.Validator = validator.New()

	isDev := cfg.AppEnv == "development"

	if isDev {
		e.Use(middleware.LoggerWithConfig(middleware.DefaultLoggerConfig))
		e.Debug = true
	}

	sessionCookieName := "SESSION_ID"

	sessionStorage := session.NewSessionStorage(mcache)

	e.Use(session.Middleware(sessionStorage, sessionCookieName, func(c echo.Context) bool {
		return strings.HasPrefix(c.Path(), "/api/v1/auth/")
	}))

	e.Use(middleware.Recover())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	healthcheckHandler := healthcheck.NewHandler()
	healthcheckHandler.RegisterRoutes(e)

	cacheStorage := cache.New(mcache)

	recipeService := recipe.NewService(logger, dbpool)
	recipeHandler := recipe.NewHandler(logger, cacheStorage, recipeService)
	recipeHandler.RegisterRoutes(e)

	passwordHasher := passwordHasher.New(&argon2id.Params{
		Memory:      cfg.ArgonMemory,
		Iterations:  cfg.ArgonIterations,
		Parallelism: cfg.ArgonParallelism,
		SaltLength:  cfg.ArgonSaltLength,
		KeyLength:   cfg.ArgonKeyLength,
	})

	userService := user.NewService(logger, passwordHasher, dbpool)

	authHandler := auth.NewHandler(logger, userService, sessionStorage, isDev, sessionCookieName)
	authHandler.RegisterRoutes(e)

	errorLog, err := zap.NewStdLogAt(logger, zapcore.ErrorLevel)
	if err != nil {
		panic(errors.Join(errors.New("failed to create a logger to http errors"), err))
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTPServerPort),
		Handler:      e,
		IdleTimeout:  time.Minute,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 8 * time.Second,
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

	_ = mcache.Close()

	fmt.Println("closed the application!")
}

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	conf "laboratory-internet-ai-test/config"
	clientgigachat "laboratory-internet-ai-test/internal/clients/gigachat"
	middlewarecors "laboratory-internet-ai-test/internal/middleware/cors"
	middlewaremetric "laboratory-internet-ai-test/internal/middleware/metric"
	middlewareratelimiter "laboratory-internet-ai-test/internal/middleware/rate_limiter"
	infrastructurepostgres "laboratory-internet-ai-test/internal/pkg/infrastructure/postgres"
	"laboratory-internet-ai-test/internal/pkg/infrastructure/redis"
	postgresrepo "laboratory-internet-ai-test/internal/repository/postgres"
	redisrepo "laboratory-internet-ai-test/internal/repository/redis"
	httptransport "laboratory-internet-ai-test/internal/transport/http"
	handlercontact "laboratory-internet-ai-test/internal/transport/http/contact_handler"
	handlermetric "laboratory-internet-ai-test/internal/transport/http/metric_handler"
	usecaselistener "laboratory-internet-ai-test/internal/usecase/contact"
	usecasemetrics "laboratory-internet-ai-test/internal/usecase/metrics"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title           Тестовое задание
// @version         1.0
// @description     Сервис контактов
// @host            localhost:3000
// @BasePath        /swagger/index.html
func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config := conf.GetConfig()
	prompts := conf.GetPrompts()

	pool, err := infrastructurepostgres.Open(ctx, config.Postgres)
	if err != nil {
		log.Panic("open postgres", "error", err)
	}
	defer pool.Close()

	redisClient, err := redis.Open(config.Redis)
	if err != nil {
		log.Panic("open postgres", "error", err)
	}

	//----------------------------

	postgresRepo := postgresrepo.New(pool, logger)
	redisRepo := redisrepo.New(redisClient, logger)

	//----------------------------

	clientGigachat := clientgigachat.New(config.GigachatAiConfig, prompts, logger)

	//----------------------------

	usecaseContact := usecaselistener.NewService(postgresRepo, clientGigachat, logger)
	usecaseMetrics := usecasemetrics.New(redisRepo, config.RateLimiter, logger)

	//----------------------------

	handlerContact := handlercontact.New(usecaseContact, logger)
	handlerMetric := handlermetric.New(usecaseMetrics, logger)

	rateLimiterMiddleware := middlewareratelimiter.NewRateLimiter(usecaseMetrics)
	statsMiddleware := middlewaremetric.NewMetric(usecaseMetrics)
	corsMiddleware := middlewarecors.CORSMiddleware()

	//----------------------------

	router := gin.Default()
	router.Use(statsMiddleware.Middleware(), rateLimiterMiddleware.Middleware(), corsMiddleware)

	httptransport.NewContactsController(router, handlerContact)
	httptransport.NewMetricController(router, handlerMetric)

	address := fmt.Sprintf(":%s", config.HTTP.Port)
	server := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	//----------------------------

	logger.Info("worker service started")

	go func() {
		logger.Info("worker service started", "host", address)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}

	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown utils server", "error", err)
	}

	logger.Info("worker service close by user")
}

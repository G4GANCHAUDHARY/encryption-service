package main

import (
	"context"
	"github.com/G4GANCHAUDHARY/encryption-service/cron"
	"github.com/G4GANCHAUDHARY/encryption-service/providers"
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener"
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/core"
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/dataMapper/dbObjectMapper"
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/dataMapper/httpDataMapper"
	"github.com/G4GANCHAUDHARY/encryption-service/urlShortener/repo"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	appConfig, err := providers.GetConfig(".env")
	if err != nil {
		panic(err)
	}

	dbClient := providers.DBClient{}
	db, err := dbClient.GetDBInstance(appConfig)
	if err != nil {
		panic(err)
	}

	urlCache := providers.GetRedisClient(appConfig, 0)
	redisCounter := providers.GetRedisClient(appConfig, 1)
	rateLimitCache := providers.GetRedisClient(appConfig, 2)

	apiRouter := mux.NewRouter()
	http.Handle("/", apiRouter)
	rateLimiter := providers.GetRateLimiter(rateLimitCache)
	handler := RateLimitMiddleware(rateLimiter, &appConfig)(apiRouter)

	urlShortenerCore := &core.UrlShortener{
		Db:               db,
		UrlRepository:    &repo.UrlRepository{},
		RedisCache:       urlCache,
		RedisCounter:     redisCounter,
		HttpResMapper:    &httpDataMapper.HttpResponseDataMapper{Config: &appConfig},
		DbMapper:         &dbObjectMapper.UrlMapper{},
		UrlAnalyticsRepo: &repo.UrlAnalyticsRepository{},
	}

	urlHandler := urlShortener.UrlHandler{
		Router:                apiRouter,
		Core:                  urlShortenerCore,
		HttpRequestDataMapper: &httpDataMapper.HttpRequestDataMapper{},
	}
	urlHandler.Init()

	server := &http.Server{
		Addr:         appConfig.HttpConfig.Address,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	urlExpiryCron := &cron.UrlExpiryCron{
		Core: urlShortenerCore,
		Db:   db,
	}

	go func() {
		log.Printf("Server starting on %s", appConfig.HttpConfig.Address)
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	/* started cron */
	go func() {
		urlExpiryCron.StartDailyCron()
	}()

	gracefulShutdown(server, &dbClient, redisCounter, urlCache, rateLimitCache, urlExpiryCron)
}

func gracefulShutdown(server *http.Server, dbClient *providers.DBClient, redisCounter *providers.RedisLib, urlCache *providers.RedisLib, rateLimitCache *providers.RedisLib, urlExpiryCron *cron.UrlExpiryCron) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Received signal: %v. Initiating graceful shutdown...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	urlExpiryCron.Stop()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	if err := dbClient.Close(); err != nil {
		log.Printf("Database close error: %v", err)
	}

	if err := redisCounter.Close(); err != nil {
		log.Printf("RedisCache counter close error: %v", err)
	}

	if err := urlCache.Close(); err != nil {
		log.Printf("url cache close error: %v", err)
	}

	if err := rateLimitCache.Close(); err != nil {
		log.Printf("rate limit close error: %v", err)
	}

	log.Println("Graceful shutdown completed")
}

func RateLimitMiddleware(rl *providers.RateLimiter, config *providers.AppConfig) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := GetIP(r)
			isAllowed := false
			if r.Method == http.MethodPost {
				isAllowed = rl.Allow(config.RateLimitConfig.WriteCapacity, r.Method+ip)
			} else {
				isAllowed = rl.Allow(config.RateLimitConfig.ReadCapacity, r.Method+ip)
			}
			if !isAllowed {
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func GetIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

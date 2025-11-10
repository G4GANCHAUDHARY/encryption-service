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

	redis := providers.GetRedisClient(appConfig)

	apiRouter := mux.NewRouter()
	http.Handle("/", apiRouter)
	urlShortenerCore := &core.UrlShortener{
		Db:               db,
		UrlRepository:    &repo.UrlRepository{},
		Redis:            redis,
		HttpResMapper:    &httpDataMapper.HttpResponseDataMapper{},
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
		Handler:      apiRouter,
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

	gracefulShutdown(server, &dbClient, redis, urlExpiryCron)
}

func gracefulShutdown(server *http.Server, dbClient *providers.DBClient, redis *providers.RedisLib, urlExpiryCron *cron.UrlExpiryCron) {
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

	if err := redis.Close(); err != nil {
		log.Printf("Redis close error: %v", err)
	}

	log.Println("Graceful shutdown completed")
}

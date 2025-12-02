package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kayaramazan/insider-message/api/cache"
	"github.com/kayaramazan/insider-message/api/database"
	"github.com/kayaramazan/insider-message/api/handler"
	"github.com/kayaramazan/insider-message/api/job"
	"github.com/kayaramazan/insider-message/api/repository"
	"github.com/kayaramazan/insider-message/api/service"
	"github.com/kayaramazan/insider-message/config"
)

func main() {

	cfg, err := config.Load("")
	if err != nil {
		log.Println("Config could not uploaded")
		return
	}

	cache, err := cache.NewRedisCache(cfg.Redis)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer cache.Close()

	db := database.NewPostgresDB(cfg.Db)
	err = db.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	messageRepo := repository.NewMessageRepository(db)

	// constructure içine al
	messageService := service.NewMessageService(messageRepo, cache, cfg.Webhook.Url, cfg.Timer.MessagePerCycle)

	job := job.New(time.Minute*time.Duration(cfg.Timer.Interval), messageService)
	go job.Start()
	handl := handler.NewHandler(messageService, job)
	SetupRoutes(handl)

}

func SetupRoutes(handler handler.Handler) {

	http.HandleFunc("PUT /api/automation/toggle", handler.StartOrStop)
	http.HandleFunc("GET /api/messages", handler.GetAllSentMessages)
	http.HandleFunc("POST /api/message", handler.CreateMessage)

	// Graceful shutdown implementation
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	log.Println("Server running on :8080")

	<-stop
	log.Println("Shutting down gracefully...")

}

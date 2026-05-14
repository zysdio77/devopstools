package main

import (
	"context"
	"log"
	"monitor_webhook/method"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := method.LoadConfig("webhook.yaml")
	if err != nil {
		log.Fatalf("load config: %s", err)
	}
	robot := method.InitRobot(cfg)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.POST("/send/:group", robot.Forward)
	r.POST("/recieve", robot.Alart)

	srv := &http.Server{
		Addr:    ":9099",
		Handler: r,
	}

	go func() {
		log.Println("server starting on :9099")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %s", err)
	}
	log.Println("server exited")
}

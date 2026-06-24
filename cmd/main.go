package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-sql-driver/mysql"
	"your-project/internal/db"
)

func main() {
	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASSWORD"),
		Net:    "tcp",
		Addr:   os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		DBName: os.Getenv("DB_NAME"),
	}

	store := db.NewMySQLStorage(cfg)
	defer func() {
		if err := store.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Your server setup here...
	// Example: srv := &http.Server{...}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Server shutting down...")
	
	// Your server shutdown here...
	// if err := srv.Shutdown(ctx); err != nil {
	//     log.Printf("Server shutdown error: %v", err)
	// }
}
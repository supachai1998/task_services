package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/supachai1998/task_services/internal/configs"
	taskRepository "github.com/supachai1998/task_services/internal/domains/tasks/infrastructure/repository"
	taskHandlerV1 "github.com/supachai1998/task_services/internal/domains/tasks/interfaces/handlers/v1"
	taskUsecase "github.com/supachai1998/task_services/internal/domains/tasks/usecases"
	"github.com/supachai1998/task_services/internal/infrastructure"
	"github.com/supachai1998/task_services/internal/interfaces"
)

// @title Task Service API
// @version 1.0
// @description This is a simple task service API.
// @termsOfService
// @contact.name Supachai
func main() {
	// Initialize configuration
	configs.InitConfig()
	// initialize database
	db, err := infrastructure.NewPostgreSQL(&configs.AppConfig.Database)
	if err != nil {
		panic("failed to connect to database")
	}
	// Initialize Echo
	e := interfaces.NewEchoInterface(&configs.AppConfig.Server)

	// Initialize repositories, use cases, and handlers
	taskRepo := taskRepository.NewTaskRepository(db)
	taskUsecase := taskUsecase.NewTaskUsecase(taskRepo)
	taskHandlerV1.NewTaskHandler(e, taskUsecase)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", configs.AppConfig.Server.Port),
		ReadTimeout:  time.Duration(configs.AppConfig.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(configs.AppConfig.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(configs.AppConfig.Server.IdleTimeout) * time.Second,
	}

	// Start the server in a goroutine
	go func() {
		if err := e.StartServer(server); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	e.Logger.Info("Gracefully shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

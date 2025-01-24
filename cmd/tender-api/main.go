package main

import (
	"Backend-trainee-assignment-autumn-2024/internal/config"
	"Backend-trainee-assignment-autumn-2024/internal/delivery/handler"
	"Backend-trainee-assignment-autumn-2024/internal/repository/postgres"
	"Backend-trainee-assignment-autumn-2024/internal/router"
	"Backend-trainee-assignment-autumn-2024/internal/service"
	"Backend-trainee-assignment-autumn-2024/internal/storage"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.NewConfig()
	db, err := storage.NewDB(cfg.DBConnStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepository := postgres.NewUserRepository(db)
	organizationRepository := postgres.NewOrganizationRepository(db)
	tenderRepository := postgres.NewTenderRepository(db)
	tenderService := service.NewTenderService(tenderRepository, userRepository, organizationRepository)
	tenderHandler := handler.NewTenderHandler(tenderService)

	app := router.SetupRouter(tenderHandler)





	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Listen(cfg.Port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	fmt.Println("Server is running on port", cfg.Port)
	<-quit
	log.Println("Shutting down server...")

	ctx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := app.ShutdownWithContext(ctx); err != nil { 
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")

}

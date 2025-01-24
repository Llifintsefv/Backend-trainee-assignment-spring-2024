package main

import (
	"Backend-trainee-assignment-autumn-2024/internal/config"
	"Backend-trainee-assignment-autumn-2024/internal/delivery/handler"
	"Backend-trainee-assignment-autumn-2024/internal/repository/postgres"
	"Backend-trainee-assignment-autumn-2024/internal/router"
	"Backend-trainee-assignment-autumn-2024/internal/service"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	cfg := config.NewConfig()
	db, err := postgres.NewDB(cfg.DBConnStr)
	if err != nil {
		slog.Error("failed to connect to database", err)
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepository := postgres.NewUserRepository(db,logger)
	organizationRepository := postgres.NewOrganizationRepository(db,logger)
	tenderRepository := postgres.NewTenderRepository(db,logger)
	tenderService := service.NewTenderService(tenderRepository, userRepository, organizationRepository,logger)
	tenderHandler := handler.NewTenderHandler(tenderService,logger)
	pingHandler := handler.NewPingHandler(logger)

	app := router.SetupRouter(tenderHandler,pingHandler)





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

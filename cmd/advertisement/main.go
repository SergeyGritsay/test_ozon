package main

import (
	"/pkg/delivery"
	"/pkg/repository"
	"/pkg/repository/postgres"
	services2 "/pkg/services"
	"advertisement-service/config"
	"log"
	"os"
	"os/signal"
	"syscall"

	"/advertisement-service/server"

	"github.com/spf13/viper"
)

func main() {
	log.Print("Starting advertisement service...")

	config.Init()

	log.Print("Connecting to the database...")
	conn, err := postgres.EstablishPSQLConnection(&postgres.PSQLConfig{
		Host:     viper.GetString("db.postgres.host"),
		Port:     viper.GetString("db.postgres.port"),
		Password: viper.GetString("db.postgres.password"),
		DBName:   viper.GetString("db.postgres.database"),
		Username: viper.GetString("db.postgres.user"),
		SSLMode:  viper.GetString("db.postgres.sslmode"),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Print("DB connection established")

	repo := repository.NewRepository(conn)
	services := services2.NewService(repo)
	handlers := delivery.NewHandler(services)

	srv := new(server.Server)
	go func() {
		if err := srv.Run(viper.GetString("app.port"), handlers.InitRoutes()); err != nil {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()
	log.Print("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("Service gracefully stopped")
}

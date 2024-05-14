package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Francesco99975/trapk/cmd/boot"
	"github.com/Francesco99975/trapk/internal/models"
)

func main() {
	err := boot.LoadEnvVariables()
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")

	models.Setup(os.Getenv("DSN"))

	e := createRouter()

	go func() {
		fmt.Printf("Running Server on port %s", port)
		e.Logger.Fatal(e.Start(":" + port))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

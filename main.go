package main

import (
	"fmt"
	"glimpseguru-tracker/router"
	"log/slog"
	"os"
)

var Logger *slog.Logger

func main() {
	r := router.New()
	errRouter := r.Run()
	if errRouter != nil {
		panic(fmt.Sprintf("Unable to run router, %e", errRouter))
	}
}

func init() {
	Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
}

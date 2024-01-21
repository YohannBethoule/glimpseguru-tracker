package main

import (
	"fmt"
	"glimpseguru-tracker/router"
	"log/slog"
	"os"
)

var logger *slog.Logger
var

func main() {
	router := router.New()
	errRouter := router.Run()
	if errRouter != nil {
		panic(fmt.Sprintf("Unable to run router, %e", errRouter))
	}
}

func init() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
}

package main

import (
	"log"
	"os"

	"github.com/rbelshevitz/gufwie/internal/appui"
	"github.com/rbelshevitz/gufwie/internal/ufw"
)

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	client := ufw.NewClient(ufw.NewExecRunner(), ufw.Config{})
	if err := appui.Run(client, logger); err != nil {
		logger.Printf("fatal: %v\n", err)
		os.Exit(1)
	}
}


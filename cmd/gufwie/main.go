package main

import (
	"flag"
	"log"
	"os"

	"github.com/rbelshevitz/gufwie/internal/appui"
	"github.com/rbelshevitz/gufwie/internal/ufw"
	"github.com/rbelshevitz/gufwie/internal/version"
)

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()
	if *showVersion {
		_, _ = os.Stdout.WriteString("gufwie " + version.Version + "\n")
		return
	}

	logger := log.New(os.Stderr, "", log.LstdFlags)
	client := ufw.NewClient(ufw.NewExecRunner(), ufw.Config{})
	if err := appui.Run(client, logger); err != nil {
		logger.Printf("fatal: %v\n", err)
		os.Exit(1)
	}
}

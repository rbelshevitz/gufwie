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
	dryRun := flag.Bool("dry-run", false, "do not make any firewall changes (shows commands instead)")
	flag.Parse()
	if *showVersion {
		line := "gufwie " + version.Version
		if version.Commit != "" && version.Commit != "dev" {
			line += " (" + version.Commit + ")"
		}
		_, _ = os.Stdout.WriteString(line + "\n")
		return
	}

	logger := log.New(os.Stderr, "", log.LstdFlags)
	client := ufw.NewClient(ufw.NewExecRunner(), ufw.Config{DryRun: *dryRun})
	if err := appui.Run(client, logger); err != nil {
		logger.Printf("fatal: %v\n", err)
		os.Exit(1)
	}
}

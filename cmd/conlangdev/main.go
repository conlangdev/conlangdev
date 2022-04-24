package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/conlangdev/conlangdev/server"
	log "github.com/sirupsen/logrus"
)

func Run() error {
	server := server.NewServer()
	if err := server.Open(); err != nil {
		return err
	}
	log.Infof("🌍 Listening on %s", server.Addr)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch

	log.Info("👋 Shutting down...")
	server.Close()
	return nil
}

func Migrate() error {
	fmt.Println("not implemented yet sorry!")
	return nil
}

func PrintUsage() {
	fmt.Println("usage: conlangdev-api [command]")
	fmt.Println("commands:")
	fmt.Println("- run: runs the web server")
	fmt.Println("- migrate: prepares sql database")
}

func main() {
	arguments := os.Args[1:]
	if len(arguments) == 0 {
		PrintUsage()
		os.Exit(1)
	}

	switch arguments[0] {
	case "run":
		if err := Run(); err != nil {
			log.WithField("command", "run").Fatal(err.Error())
		}
		os.Exit(0)
	case "migrate":
		if err := Migrate(); err != nil {
			log.WithField("command", "migrate").Fatal(err.Error())
		}
		os.Exit(0)
	default:
		PrintUsage()
		os.Exit(1)
	}
}
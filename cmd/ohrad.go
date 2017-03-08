package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bertinatto/ohrad"
)

func handleSignals(cleanup func() error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGCHLD, syscall.SIGHUP)
	go func() {
		for s := range c {
			switch s {
			case os.Interrupt, syscall.SIGTERM:
				log.Println("Gracefully shutting down the daemon")
				if err := cleanup(); err != nil {
					log.Fatal(err)
					os.Exit(-1)
				} else {
					log.Println("Shut down finished")
					os.Exit(0)
				}
			}
		}
	}()
}

func main() {
	server := ohrad.Server{
		Addr:        ":123",
		ReadTimeout: 5,
	}
	handleSignals(server.Shutdown)
	log.Fatal(server.ListenAndServe())
}

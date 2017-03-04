package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bertinatto/ohrad"
)

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGCHLD, syscall.SIGHUP)

}

func main() {
	server := ohrad.Server{
		Addr:        ":123",
		ReadTimeout: 5,
	}
	log.Fatal(server.ListenAndServe())
}

/*
 * TODO:
 * signal handling
 * privilege separation
 * logging
 * handler
 * confs
 */

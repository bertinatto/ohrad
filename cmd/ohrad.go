package main

import (
	"os"

	"github.com/bertinatto/ohrad"
)

func main() {
	server := ohrad.Server{
		Addr:        ":123",
		ReadTimeout: 5,
	}
	ohrad.HandleSignals(server.Shutdown)
	if err := server.ListenAndServe(); err != nil {
		ohrad.Log.Err(err.Error())
		os.Exit(-1)
	}
	os.Exit(0)
}

package ohrad

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

func unix2Ntp(tsp int64) uint32 {
	return uint32(Jan1970 + tsp)
}

func unix2NtpNanos(tsp int64) uint32 {

	return uint32((Jan1970 * 1000 * 1000 * 1000) + tsp)
}

func ntp2Unix(tsp int64) uint32 {
	return uint32(tsp - Jan1970)
}

func HandleSignals(cleanup func() error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGCHLD, syscall.SIGHUP)
	go func() {
		for s := range c {
			switch s {
			case os.Interrupt, syscall.SIGTERM:
				Log.Info("Gracefully shutting down the daemon")
				if err := cleanup(); err != nil {
					Log.Err(err.Error())
					os.Exit(-1)
				} else {
					Log.Info("Shut down finished")
					os.Exit(0)
				}
			}
		}
	}()
}

func getTimeNow() NtpLong {
	t := time.Now()
	secs := t.Unix()
	nanos := t.UnixNano()
	left := nanos - (secs * NanosPerSecond)
	fraction := float32(left) / float32(NanosPerSecond)
	return NtpLong{
		IntParl:   unix2Ntp(secs),
		Fractionl: uint32(float32(fraction) * float32(NanosPerSecond)),
	}
}

package ohrad

import (
	"fmt"
	"log/syslog"
)

type Logger struct {
	*syslog.Writer
}

var Log = new(Logger)

func init() {
	var err error
	Log.Writer, err = syslog.New(syslog.LOG_INFO, "ohrad")
	if err != nil {
		panic("Couldn't create Syslog writer")
	}
}

func (w *Logger) NtpMsg(msg *NtpMsg) {
	w.Debug(fmt.Sprintf("NTP message: %+v", *msg))
}

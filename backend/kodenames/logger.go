package kodenames

import (
	"os"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("kodenames")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s}%{color:reset} %{message}`,
)

func SetupLogging() {
	// For demo purposes, create two backend for os.Stderr.
	backend1 := logging.NewLogBackend(os.Stdout, "", 0)

	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backend2Formatter := logging.NewBackendFormatter(backend1, format)


	// Set the backends to be used.
	logging.SetBackend(backend2Formatter)
}
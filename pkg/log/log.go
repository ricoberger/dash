package log

import (
	"log"
	"os"
	"path/filepath"
)

var (
	f     *os.File
	debug bool
)

func Init(dir string, debugEnabled bool) error {
	debug = debugEnabled

	var err error

	dashLog := filepath.Join(dir, "dash.log")

	f, err = os.OpenFile(dashLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	log.SetOutput(f)
	return nil
}

func Close() {
	f.Close()
}

func Debugf(format string, v ...interface{}) {
	if debug {
		log.Printf(format, v...)
	}
}

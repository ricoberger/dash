package log

import (
	"log"
	"os"
)

var (
	f     *os.File
	debug bool
)

func Init(dir string, debugEnabled bool) error {
	debug = debugEnabled

	var err error

	f, err = os.OpenFile(dir+"/dash.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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

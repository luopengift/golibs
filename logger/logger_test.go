package logger

import (
	"testing"
)

func Test_logger(t *testing.T) {
	MyLogger.SetLevel(TRACE)
	Trace("1234567")
	Debug("<%s,%s>", "hello", "xxx")
	Info("hello")
	Warn("hello", "xxx")
	Error("hello", "xxx")
	Fatal("hello", "xxx")
	//Panic("hello", "xxx")

}

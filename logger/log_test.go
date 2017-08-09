package logger

import (
	"testing"
)

func Test_logger(t *testing.T) {
	MyLogger.SetLevel(TRACE)
	MyLogger.SetPrefix("Module")
	Trace("1234567")
	Debug("<%s,%s>", "hello", "xxx")
	Info("hello")
	Warn("hello", "xxx")
	Error("hello", "xxx")
	Fatal("hello", "xxx")
	//Panic("hello", "xxx")

}

func Benchmark_logger(b *testing.B) {
        MyLogger.SetLevel(PANIC)
        MyLogger.SetPrefix("Module")
        for i := 0; i < b.N; i++ {
                Trace("1234567")
                Debug("<%s,%s>", "hello", "xxx")
                Info("hello")
                Warn("hello", "xxx")
                Error("hello", "xxx")
                Fatal("hello", "xxx")
                //Panic("hello", "xxx")
        }
}

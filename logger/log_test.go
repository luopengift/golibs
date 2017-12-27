package logger

import (
	"os"
	"testing"
)

func Test_FileHandler(t *testing.T) {
	f := NewFileWriter("tmp/test_%Y%M%D%h%m%s.log", 1)
	MyLogger.SetOutput(f, os.Stdout)
	Info("@@@@@@@@@@@")
}

func Test_logger(t *testing.T) {
	MyLogger.SetLevel(NULL)
	MyLogger.SetPrefix("Module")
	Trace("1234567")
	Debug("<%s,%s>", "hello", "xxx")
	Info("hello")
	Warn("hello %s", "xxx")
	Error("hello %s", "xxx")
	Fatal("hello %s", "xxx")
	//Panic("hello", "xxx")

}

func Benchmark_logger(b *testing.B) {
	MyLogger.SetLevel(PANIC)
	MyLogger.SetPrefix("Module")
	for i := 0; i < b.N; i++ {
		Trace("1234567")
		Debug("<%s,%s>", "hello", "xxx")
		Info("hello")
		Warn("hello %s", "xxx")
		Error("hello %s", "xxx")
		Fatal("hello %s", "xxx")
		//Panic("hello", "xxx")
	}
}

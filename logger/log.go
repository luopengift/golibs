package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

const (
	TRACE   = uint8(1 << iota) //1
	DEBUG                      //2
	INFO                       //4
	WARNING                    //8
	ERROR                      //16
	FATAL                      //32
	PANIC                      //32
	OFF                        //64
)

var (
	Files = [3]*os.File{os.Stdin, os.Stdout, os.Stderr}
	Level = map[uint8]string{
		TRACE:   "[T]",
		DEBUG:   "[D]",
		INFO:    "[I]",
		WARNING: "[W]",
		ERROR:   "[E]",
		FATAL:   "[F]",
		PANIC:   "[P]",
	}
)

type Logger struct {
	lv     uint8
	prefix string
	color  bool
	out    []io.Writer
}

func NewLogger(lv uint8, prefix string, color bool, out ...io.Writer) *Logger {
	return &Logger{
		lv:     lv,
		prefix: prefix,
		color:  color,
		out:    out,
	}
}

func (self *Logger) SetLevel(lv uint8) {
	self.lv = lv
}

func (self *Logger) SetPrefix(prefix string) {
	self.prefix = prefix
}

func (self *Logger) SetColor(color bool) {
	self.color = color
}

func (self *Logger) SetOutput(out ...io.Writer) {
	self.out = out
}

func (self *Logger) format(lv uint8, format string) string {
	str := ""
	if self.prefix != "" {
		str += fmt.Sprintf("%s %s ", time.Now().Format(self.prefix), Level[lv])
	}
	str += format
	if self.color {
		str = setColor(lv, str)
	}
	return str
}

func (self *Logger) writeLog(lv uint8, format string, msg ...interface{}) error {
	if lv < self.lv {
		return nil
	}
	self.output(self.format(lv, format), msg...)
	return nil
}

func (self *Logger) output(format string, msg ...interface{}) {
	for _, out := range self.out {
		fmt.Fprintf(out, format+"\n", msg...)
	}
}

func (self *Logger) Trace(format string, msg ...interface{}) {
	self.writeLog(TRACE, format, msg...)
}

func (self *Logger) Debug(format string, msg ...interface{}) {
	self.writeLog(DEBUG, format, msg...)
}
func (self *Logger) Info(format string, msg ...interface{}) {
	self.writeLog(INFO, format, msg...)
}
func (self *Logger) Warn(format string, msg ...interface{}) {
	self.writeLog(WARNING, format, msg...)
}
func (self *Logger) Error(format string, msg ...interface{}) {
	self.writeLog(ERROR, format, msg...)
}
func (self *Logger) Fatal(format string, msg ...interface{}) {
	self.writeLog(FATAL, format, msg...)
}
func (self *Logger) Panic(format string, msg ...interface{}) {
	self.writeLog(PANIC, format, msg...)
	panic(fmt.Sprintf(format, msg...))
}

package logger

import (
	"bytes"
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
	PANIC                      //64
	OFF                        //128
	NULL    = uint8(0)
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
	lv         uint8  //最低日志级别
	timeFormat string //时间格式
	prefix     string
	color      bool
	out        []io.Writer
}

func NewLogger(timeFormat string, lv uint8, out ...io.Writer) *Logger {
	return &Logger{
		lv:         lv,
		timeFormat: timeFormat,
		prefix:     "",
		color:      true,
		out:        out,
	}
}

func (log *Logger) SetLevel(lv uint8) {
	log.lv = lv
}

func (log *Logger) SetTimeFormat(timeFormat string) {
	log.timeFormat = timeFormat
}

func (log *Logger) SetPrefix(prefix string) {
	log.prefix = prefix
}

func (log *Logger) SetColor(color bool) {
	log.color = color
}

func (log *Logger) SetOutput(out ...io.Writer) {
	log.out = out
}

func (log *Logger) format(lv uint8, format string) string {
	var buf bytes.Buffer
	if log.timeFormat != "" {
		buf.WriteString(time.Now().Format(log.timeFormat))
		buf.WriteString(" ")
	}
	if log.lv != NULL {
		buf.WriteString(Level[lv])
		buf.WriteString(" ")
	}
	if log.prefix != "" {
		buf.WriteString(log.prefix)
		buf.WriteString(" ")
	}
	buf.WriteString(format)
	buf.WriteString(" ")
	str := buf.String()
	if log.color {
		str = setColor(lv, str)
	}
	return str
}

func (log *Logger) writeLog(lv uint8, format string, msg ...interface{}) error {
	if lv < log.lv {
		return nil
	}
	log.output(log.format(lv, format), msg...)
	return nil
}

func (log *Logger) output(format string, msg ...interface{}) {
	for _, out := range log.out {
		fmt.Fprintf(out, format+"\n", msg...)
	}
}

func (log *Logger) Trace(format string, msg ...interface{}) {
	log.writeLog(TRACE, format, msg...)
}

func (log *Logger) Debug(format string, msg ...interface{}) {
	log.writeLog(DEBUG, format, msg...)
}
func (log *Logger) Info(format string, msg ...interface{}) {
	log.writeLog(INFO, format, msg...)
}
func (log *Logger) Warn(format string, msg ...interface{}) {
	log.writeLog(WARNING, format, msg...)
}
func (log *Logger) Error(format string, msg ...interface{}) {
	log.writeLog(ERROR, format, msg...)
}
func (log *Logger) Fatal(format string, msg ...interface{}) {
	log.writeLog(FATAL, format, msg...)
}
func (log *Logger) Panic(format string, msg ...interface{}) {
	log.writeLog(PANIC, format, msg...)
	panic(fmt.Sprintf(format, msg...))
}

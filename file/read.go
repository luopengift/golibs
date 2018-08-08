package file

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/luopengift/golibs/logger"
)

type Tail struct {
	*File
	cname    string //config name
	line     chan []byte
	reader   *bufio.Reader
	interval int64
	Handler  // file name handle interface
	//EOF
	//@ true: stop
	//@ false: wait
	endstop bool
}

func NewTail(cname string, handler Handler) *Tail {
	name := handler.Handle(cname)
	file := NewFile(name, os.O_RDONLY)

	tail := new(Tail)
	tail.File = NewFile(name, os.O_RDONLY)
	tail.cname = cname
	tail.line = make(chan []byte, 1)
	tail.reader = bufio.NewReader(file.fd)
	tail.interval = 1000 //ms
	tail.endstop = false
	tail.Handler = handler
	return tail
}

func (self *Tail) EndStop(b bool) {
	self.endstop = b
}

func (self *Tail) Close() error {
	close(self.line)
	return self.File.Close()
}

func (self *Tail) ReOpen() error {
	if err := self.File.Close(); err != nil {
		logger.Error("<file %v close fail:%v>", self.name, err)
		return err
	}
	self.name = self.Handler.Handle(self.cname)
	err := self.Open()
	if err != nil {
		return err
	}
	self.reader = bufio.NewReader(self.fd)
	return nil
}

func (self *Tail) Stop() {
	self.File.Close()
	close(self.line)
}

func (self *Tail) ReadLine() {
	go func() {

		offset, err := self.TrancateOffsetByLF(self.seek)
		if err != nil {
			logger.Error("<Trancate offset:%d,Error:%+v>", self.seek, err)
		}
		err = self.Seek(offset)
		if err != nil {
			logger.Error("<seek offset[%d] error:%+v>", self.seek, err)
		}

		for {
			line, err := self.reader.ReadBytes('\n')
			switch {
			case err == io.EOF:
				if self.endstop {
					logger.Info("<file %s is END:%+v>", self.name, err)
					self.Stop()
					return
				}
				time.Sleep(time.Duration(self.interval) * time.Millisecond)
				if self.name == self.cname {
					if self.IsSameFile(self.name) {
						continue
					} else {
						self.ReOpen()
					}
				} else {
					if self.name == self.Handler.Handle(self.cname) { //检测是否需要按时间轮转新文件
						continue
					} else {
						self.ReOpen()
					}
				}

			case err != nil && err != io.EOF:
				time.Sleep(time.Duration(self.interval) * time.Millisecond)
				logger.Error("<Read file error:%v,%v>", line, err)
				self.ReOpen()
				continue
			default:
				msg := bytes.TrimRight(line, "\n")
				self.line <- msg
				self.seek += int64(len(line))
			}
		}
	}()
}

func (self *Tail) NextLine() chan []byte {
	return self.line
}

func (self *Tail) Read(p []byte) (int, error) {
	msg, ok := <-self.line
	if !ok {
		return 0, fmt.Errorf("file is closed")
	}
	if len(msg) > len(p) {
		return 0, errors.New("message is large than buf")
	}
	n := copy(p, msg)
	return n, nil
}

package tail

import (
	"bufio"
	"github.com/luopengift/golibs/logger"
	"io"
	"os"
	"strings"
	"time"
)

type Tail struct {
	filename string
	line     chan *string
	file     *os.File
	reader   *bufio.Reader
	interval int64
}

func NewTail(filename string) *Tail {
	file, err := os.Open(filename)
	if err != nil {
		logger.Panic("<%s can not open>:%v", filename, err)
		return nil
	}
	return &Tail{
		filename: filename,
		line:     make(chan *string),
		file:     file,
		reader:   bufio.NewReader(file),
		interval: 500, //ms
	}
}

func (self *Tail) ReadLine() {
	go func() {
		for {
			line, err := self.reader.ReadString('\n')
			switch {
			case err == io.EOF:
				time.Sleep(time.Duration(self.interval) * time.Millisecond)
			case err != nil && io.EOF != err:
				continue
			default:
				msg := strings.TrimRight(line, "\n")
				self.line <- &msg
			}
		}
	}()
}

func (self *Tail) NextLine() chan *string {
	return self.line
}

func (self *Tail) Stop() {
	self.file.Close()
	close(self.line)
}

func (self *Tail) Offset() int64 {
	offset, _ := self.file.Seek(0, os.SEEK_CUR)
	return offset
}

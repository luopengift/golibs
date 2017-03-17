package file

import (
	"bufio"
	"github.com/luopengift/golibs/logger"
	"io"
	"os"
	"strings"
	"time"
)

type Tail struct {
	*File
	cname    string //config name
	line     chan *string
	reader   *bufio.Reader
	interval int64
}

func NewTail(cname string) *Tail {
	name := HandlerRule(cname)
	file := NewFile(name, os.O_RDONLY)
	return &Tail{
		file,
		cname,
		make(chan *string),
		bufio.NewReader(file.fd),
		1000, //ms
	}
}

func (self *Tail) ReOpen() {
	if err := self.Close(); err != nil {
		logger.Error("<file %v close fail:%v>", self.name, err)
	}
	self.name = HandlerRule(self.cname)
	err := self.Open()
	if err != nil {
		return
	}
	self.reader = bufio.NewReader(self.fd)
}

func (self *Tail) Stop() {
	self.Close()
	close(self.line)
}

func (self *Tail) ReadLine() {
	go func() {
		for {
			line, err := self.reader.ReadString('\n')
			switch {
			case err == io.EOF:
				time.Sleep(time.Duration(self.interval) * time.Millisecond)
				if self.name == self.cname {
					if inode, err := Inode(self.name); err != nil { //检测是否需要重新打开新的文件
						continue
					} else {
						if inode != self.inode {
							self.ReOpen()
						}
					}
				} else {
					if self.name == HandlerRule(self.cname) { //检测是否需要按时间轮转新文件
						continue
					} else {
						self.ReOpen()
					}
				}

			case err != nil && err != io.EOF:
				logger.Error("<Read file error:%v,%v>", line, err)
				self.ReOpen()
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

func (self *Tail) Offset() int64 {
	offset, _ := self.fd.Seek(0, os.SEEK_CUR)
	return offset
}

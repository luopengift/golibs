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
	line     chan *string
	reader   *bufio.Reader
	interval int64
}

func NewTail(name string) *Tail {
	file := NewFile(name)
	return &Tail{
		file,
		make(chan *string),
		bufio.NewReader(file.fd),
		1000, //ms
	}
}

func (self *Tail) ReOpen() {
	if err := self.Close(); err != nil {
		logger.Error("<file %v close fail:%v>", self.name, err)
		return
	}
	fd, err := self.Open()
	if err != nil {
		return
	}
	self.fd = fd
	inode, err := self.Inode()
	if err != nil {
		return
	}
	self.inode = inode
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
				if inode, err := Inode(self.name); err != nil { //检测是否需要重新打开新的文件
					logger.Debug("file %s get inode error:%v", self.name, err)
					continue
				} else {
					logger.Debug("file get inode success:%v", inode)
					if inode != self.inode {
						logger.Debug("inode is not same%v,%v", inode, self.inode)
						self.ReOpen()
					}
				}
				time.Sleep(time.Duration(self.interval) * time.Millisecond)
			case err != nil && err != io.EOF:
				logger.Error("<Read file error:%v,%v>", line, err)
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

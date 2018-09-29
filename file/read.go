package file

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/luopengift/log"
)

// Tail tail
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

// NewTail new tail
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

// EndStop end stop
func (t *Tail) EndStop(b bool) {
	t.endstop = b
}

// Close close
func (t *Tail) Close() error {
	close(t.line)
	return t.File.Close()
}

// ReOpen re open
func (t *Tail) ReOpen() error {
	if err := t.File.Close(); err != nil {
		log.Error("<file %v close fail:%v>", t.name, err)
		return err
	}
	t.name = t.Handler.Handle(t.cname)
	err := t.Open()
	if err != nil {
		return err
	}
	t.reader = bufio.NewReader(t.fd)
	return nil
}

// Stop stop
func (t *Tail) Stop() {
	t.File.Close()
	close(t.line)
}

// ReadLine read line
func (t *Tail) ReadLine() {
	go func() {

		offset, err := t.TrancateOffsetByLF(t.seek)
		if err != nil {
			log.Error("<Trancate offset:%d,Error:%+v>", t.seek, err)
		}
		err = t.Seek(offset)
		if err != nil {
			log.Error("<seek offset[%d] error:%+v>", t.seek, err)
		}

		for {
			line, err := t.reader.ReadBytes('\n')
			switch {
			case err == io.EOF:
				if t.endstop {
					log.Info("<file %s is END:%+v>", t.name, err)
					t.Stop()
					return
				}
				time.Sleep(time.Duration(t.interval) * time.Millisecond)
				if t.name == t.cname {
					if t.IsSameFile(t.name) {
						continue
					} else {
						t.ReOpen()
					}
				} else {
					if t.name == t.Handler.Handle(t.cname) { //检测是否需要按时间轮转新文件
						continue
					} else {
						t.ReOpen()
					}
				}

			case err != nil && err != io.EOF:
				time.Sleep(time.Duration(t.interval) * time.Millisecond)
				log.Error("<Read file error:%v,%v>", line, err)
				t.ReOpen()
				continue
			default:
				msg := bytes.TrimRight(line, "\n")
				t.line <- msg
				t.seek += int64(len(line))
			}
		}
	}()
}

// NextLine nextline
func (t *Tail) NextLine() chan []byte {
	return t.line
}

// Read read
func (t *Tail) Read(p []byte) (int, error) {
	msg, ok := <-t.line
	if !ok {
		return 0, fmt.Errorf("file is closed")
	}
	if len(msg) > len(p) {
		return 0, errors.New("message is large than buf")
	}
	n := copy(p, msg)
	return n, nil
}

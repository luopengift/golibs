package file

import (
	"github.com/luopengift/golibs/logger"
	"io/ioutil"
	"os"
	"syscall"
)

type File struct {
	name  string
	model int
	inode uint64
	fd    *os.File
	seek  int64
}

func NewFile(name string, model int) *File {
	file := &File{name: name, model: model}
	file.Open()
	return file
}

func (self *File) Open() (err error) {
	self.fd, err = os.OpenFile(self.name, self.model, 0660)
	if err != nil {
		logger.Error("<file %s can not open>:%v", self.name, err)
		return err
	}

	self.inode, err = self.Inode()
	if err != nil {
		logger.Error("< %s can not get inode>:%v", self.name, err)
		return err
	}
	return nil
}

func (self *File) Close() error {
	return self.fd.Close()
}

func (self *File) Fd() *os.File {
	return self.fd
}

func (self *File) ReadAll() (file []byte, err error) {
	file, err = ioutil.ReadAll(self.fd)
	return
}

func Inode(name string) (uint64, error) {
	if stat, err := os.Stat(name); err != nil {
		return 0, err
	} else {
		inode := stat.Sys().(*syscall.Stat_t).Ino
		return inode, nil
	}
}

func (self *File) Inode() (uint64, error) {
	inode, err := Inode(self.name)
	return inode, err
}

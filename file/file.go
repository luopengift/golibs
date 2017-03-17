package file

import (
	"github.com/luopengift/golibs/logger"
	"os"
	"syscall"
)

type File struct {
	name  string
	inode uint64
	fd    *os.File
	seek  int64
}

func NewFile(name string) *File {
	file := &File{}

	file.name = name

	fd, err := file.Open()
	if err != nil {
		logger.Error("<file %s can not open>:%v", name, err)
		return nil
	}
	file.fd = fd

	inode, err := file.Inode()
	if err != nil {
		logger.Error("< %s can not get inode>:%v", name, err)
		return nil
	}
	file.inode = inode

	return file
}

func (self *File) Open() (*os.File, error) {
	return os.OpenFile(self.name, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
}

func (self *File) Close() error {
	return self.fd.Close()
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

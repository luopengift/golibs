package file

import (
	"github.com/luopengift/golibs/logger"
	"io/ioutil"
	"os"
	"path/filepath"
)

type File struct {
	name  string
	model int
	fd    *os.File
	seek  int64
}

func NewFile(name string, model int) *File {
	file := &File{
		name:  name,
		model: model,
	}
	file.Open()
	return file
}

// 文件名称
func (f *File) Name() string {
	return f.name
}

func (f *File) Dir() string {
	return filepath.Dir(f.name)
}

// 文件名
func (f *File) BaseName() string {
	return filepath.Base(f.name)
}

func (f *File) Open() (err error) {
	f.fd, err = os.OpenFile(f.name, f.model, 0660)
	if err != nil {
		logger.Error("<file %s can not open>:%v", f.name, err)
		return
	}
	return
}

func (f *File) Read(p []byte) (int, error) {
	return f.fd.Read(p)
}

func (f *File) Write(p []byte) (int, error) {
	return f.fd.Write(p)
}

func (f *File) Close() error {
	return f.fd.Close()
}

func (f *File) Fd() *os.File {
	return f.fd
}

// os.SEEK_CUR int = 1 // seek relative to the current offset
// os.SEEK_SET int = 0 // seek relative to the origin of the file
// os.SEEK_END int = 2 // seek relative to the end
func (f *File) Seek(offset int64) (err error) {
	f.seek, err = f.fd.Seek(offset, os.SEEK_SET)
	return
}

func (f *File) ReadOneByte(offset int64) ([]byte, error) {
	buf := make([]byte, 1)
	_, err := f.fd.ReadAt(buf, offset)
	return buf, err

}

// 根据offset值,往前计算该行的起始偏移量
func (f *File) TrancateOffsetByLF(offset int64) (int64, error) {
	for ; offset >= 0; offset-- {
		buf, err := f.ReadOneByte(offset)
		if err != nil {
			return 0, err
		}
		if string(buf) == "\n" {
			return offset + 1, nil //pos为"\n"的位置,需要加1才是行首的位置
		}
	}
	return 0, nil
}

// 根据offset值,往后计算该行的起始偏移量
func (f *File) CeilingOffsetByLF(offset int64) (int64, error) {
	for ; ; offset++ {
		buf, err := f.ReadOneByte(offset)
		if err != nil {
			return 0, err
		}
		if string(buf) == "\n" {
			return offset + 1, nil //pos为"\n"的位置,需要加1才是行首的位置
		}
	}
	return 0, nil
}

func (f *File) Offset() int64 {
	return f.seek
}

func (f *File) Size() int64 {
	if stat, err := f.fd.Stat(); err != nil {
		return 0
	} else {
		return stat.Size()
	}
}

func (f *File) ReadAll() (file []byte, err error) {
	file, err = ioutil.ReadAll(f.fd)
	return
}

func (f *File) IsSameFile(file string) bool {
    stat1, err := os.Stat(file)
    if err != nil {
        logger.Warn("SameFile error:%v",err)
        return false
    }
    stat2, err := f.fd.Stat()
    if err != nil {
        logger.Warn("SameFile error:%v",err)
        return false
    }
    return os.SameFile(stat1, stat2)

}

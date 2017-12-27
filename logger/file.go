package logger

import (
	"github.com/luopengift/types"
	"io"
	"io/ioutil"
	"path"
	"os"
	"strconv"
)

type FileWriter struct {
	Name     string //标记名称
	timeName string //时间文件名
	RealName string //实际名称
	MaxSize  int64  //最大文件大小
	count    int64    //轮转文件计数
}

func NewFileWriter(name string, max int64) io.Writer {
	return &FileWriter{
		Name:     name,
		timeName: types.TimeFormat(name),
		RealName: types.TimeFormat(name),
		MaxSize:  max,
	}
}

func (f *FileWriter) handler() string {
	err := os.MkdirAll(path.Dir(f.RealName), 0755)
	if err != nil {
		println(err.Error())
		return f.RealName
	}
	info, err := os.Stat(f.RealName)
	if err != nil {
		println(err.Error())
		return f.RealName
	}
	timeName := types.TimeFormat(f.Name)
	if timeName != f.timeName {
		f.timeName = timeName
		f.count = 0
	}
	if info.Size() > f.MaxSize {
		f.count += 1
		f.RealName = f.timeName + "." + strconv.FormatInt(f.count, 10)
	}
	return f.RealName
}

func (f *FileWriter) Write(p []byte) (int, error) {
	f.handler()
	if err := ioutil.WriteFile(f.RealName, p, 0644); err != nil {
		return 0, err
	}
	return len(p), nil
}



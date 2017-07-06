package file

import (
	"regexp"
	"time"
)

var TimeRule *TimeRuler
var NullRule *NullRuler

var (
	//时间通配符，用于正则表达式替换
	Map map[string]string = map[string]string{
		"%Y": "2006",
		"%M": "01",
		"%D": "02",
		"%h": "15",
		"%m": "04",
		"%s": "05",
	}
)

// 处理接口
type Handler interface {
	Handle(string) string
}

type TimeRuler struct{}

//eg:"test-%Y%M%D.log" ->"test-20170203.log"
//eg:"test-%Y-%M-%D.log" ->"test-2017-02-03.log"
func (t *TimeRuler) Handle(str string) string {
	for k, v := range Map {
		re, err := regexp.Compile(k)
		if err != nil {
			continue
		}
		str = re.ReplaceAllString(str, v)
	}
	return time.Now().Format(str)
}

type NullRuler struct{}

func (n *NullRuler) Handle(str string) string {
	return str
}

func init() {
	TimeRule = new(TimeRuler)
	NullRule = new(NullRuler)
}

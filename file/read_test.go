package file

import (
	"fmt"
	"testing"
)

func Test_read(t *testing.T) {
	//tt := NewTail("test-%Y-%M-%D-%h-%m.log")
	tt := NewTail("test.log", &TimeRuler{})
	tt.ReadLine()
	tt.EndStop(true)
	for v := range tt.NextLine() {
		fmt.Println(string(v))
	}
}

type RuntimeConfig struct {
	DEBUG    bool `json:"DEBUG"`
	MAXPROCS int  `json:"MAXPROCS"`
}

type KafkaConfig struct {
	Addrs      []string `json:"addrs"`
	Topic      string   `json:"topic"`
	MaxThreads int64    `json:"maxthreads"`
}

type HttpConfig struct {
	Addr string `json:"addr"`
}

type TestConfig struct {
	Runtime RuntimeConfig
	Kafka   KafkaConfig
	File    []string `json:"file"`
	Prefix  string   `json:"prefix"`
	Suffix  string   `json:"suffix"`
	Http    HttpConfig
	Tags    string `json:"tags"`
	Version string `json:version`
}

func Test_config(t *testing.T) {
	test := &TestConfig{}
	config := NewConfig("./config.json")
	config.Parse(test)
	fmt.Println(fmt.Sprintf("%+v", test))
	fmt.Println(config)

}


func Test_tail(t *testing.T) {
    f := NewTail("/tmp/test/t_%Y%M%D%h%m.log",&TimeRuler{})
    f.ReadLine()
    f.EndStop(false)
    for v := range f.NextLine() {
        fmt.Println(string(v))
    }
}

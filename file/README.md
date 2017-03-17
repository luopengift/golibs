## 读取一个文件

#### 类似于shell的tail -F,不过只能从文件开始位置读起
#### 可以按照时间匹配文件(%Y:年 %M:月 %D:日期 %h:小时 %m:分钟)
#### TODO：增加读取的起始位置

### 使用例子:
```
package main

import (
    "github.com/luopengift/golibs/file"
    "fmt"
)

func main() {
    f := file.NewTail("test-%Y-%M-%D.log")
    f.ReadLine()

    for v := range f.NextLine() {
        fmt.Println(*v) //NextLine返回行内容的地址
    }

    f.Stop()
}

```

## 读取一个文件

#### 类似于shell的tail -F,不过只能从文件开始位置读起
#### TODO：增加读取的起始位置

### 使用例子:
```
package main

import (
    "github.com/luopengift/golibs/file"
    "fmt"
)

func main() {
    f := file.NewTail("test.log")
    f.ReadLine()

    for v := range f.NextLine() {
        fmt.Println(*v) //NextLine返回行内容的地址
    }

    f.Stop()
}

```

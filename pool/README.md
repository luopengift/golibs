## 控制golang协程并发数量

### 使用例子:
```
package main

import (
    "github.com/luopengift/golibs/pool"
    "time"
    "fmt"
)

func main() {
    var max_works int64 = 100
    p := pool.NewPool(max_works,nil)
    for i:=0;i<1000;i++ {
        go p.Run(func() error {
            fmt.Println(p)
            time.Sleep(1*time.Second)
            return nil
        })
    }
    p.Wait()
}

```

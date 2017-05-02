## golang 队列

### 使用例子:
```
package main

import (
    "github.com/luopengift/golibs/channel"
    "time"
    "fmt"
)

func main() {
    var max_works int64 = 100
    _ := channel.NewChannel(max_works)
}

```

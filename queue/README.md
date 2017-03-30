## golang 队列

### 使用例子:
```
package main

import (
    "github.com/luopengift/golibs/queue"
    "time"
    "fmt"
)

func main() {
    var max_works int64 = 100
    _ := queue.NewQueue(max_works,nil)
}

```

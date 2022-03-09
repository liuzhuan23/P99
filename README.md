# P99
go版本的p99统计调用，区别与C版本的，这个版本直接import即可

#### go版本p99
```
同样摘自rpcx原作者smallnest的benchmark，附原文连接

https://github.com/rpcxio/rpcx-benchmark.git
```

#### usage example
```
package main

import "abc/P99"
import "fmt"

type aaa struct {
	id int
}

func (a *aaa) Init() int {
    return 1
}

func (a *aaa) Run() int {
	fmt.Println("aaa id: ", a.id)
    return 1
}

type bbb struct {
	id int
}

func (b *bbb) Init() int {
    return 1
}

func (b *bbb) Run() int {
	fmt.Println("bbb id: ", b.id)
    return 1
}

func main() {
    f1 := P99.NewP99Stat(8, 4)

    var a1 aaa
    a1.id = 1
    f1.P99Register(a1.Init, a1.Run)

    var a2 aaa
    a2.id = 2
    f1.P99Register(a2.Init, a2.Run)

    var b1 bbb
    b1.id = 3
    f1.P99Register(b1.Init, b1.Run)

    var b2 bbb
    b2.id = 4
    f1.P99Register(b2.Init, b2.Run)

    fmt.Println("数量: ", f1.P99IfCount())

    f1.P99Run()
}
```

#### example output
```
数量:  4
bbb id:  4
bbb id:  4
aaa id:  1
aaa id:  1
aaa id:  2
aaa id:  2
bbb id:  3
bbb id:  3
2022/03/09 13:46:18 P99.go:106: INFO : P99客户端数量: 4
2022/03/09 13:46:18 P99.go:107: INFO : P99每个客户端的并发请求数: 2
2022/03/09 13:46:18 stats.go:15: INFO : took 0 ms for 8 requests
2022/03/09 13:46:18 stats.go:36: INFO : sent     requests    : 8
2022/03/09 13:46:18 stats.go:37: INFO : received requests    : 8
2022/03/09 13:46:18 stats.go:38: INFO : received requests_OK : 0
2022/03/09 13:46:18 stats.go:40: INFO : throughput  (TPS)    : 148961
2022/03/09 13:46:18 stats.go:45: INFO : mean: 822 ns, median: 744 ns, max: 1252 ns, min: 696 ns, p99.9: 1096 ns
2022/03/09 13:46:18 stats.go:46: INFO : mean: 0 ms, median: 0 ms, max: 0 ms, min: 0 ms, p99.9: 0 ms
```

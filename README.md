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

import (
	"github.com/liuzhuan23/P99"
)

func aaa() int {
	fmt.Println("test123")
	return 12
}

func main() {
	f1 := P99.NewP99Stat(10, 1)
	f1.Register(aaa)
	f1.Run()
}
```

#### example output
```
test123
test123
test123
test123
test123
test123
test123
test123
test123
test123
2022/03/08 11:55:31 P99.go:86: INFO : P99客户端数量: 1
2022/03/08 11:55:31 P99.go:87: INFO : P99每个客户端的并发请求数: 10
2022/03/08 11:55:31 stats.go:15: INFO : took 0 ms for 10 requests
2022/03/08 11:55:31 stats.go:36: INFO : sent     requests    : 10
2022/03/08 11:55:31 stats.go:37: INFO : received requests    : 10
2022/03/08 11:55:31 stats.go:38: INFO : received requests_OK : 10
2022/03/08 11:55:31 stats.go:40: INFO : throughput  (TPS)    : 32467
2022/03/08 11:55:31 stats.go:45: INFO : mean: 14470 ns, median: 10941 ns, max: 45679 ns, min: 10001 ns, p99.9: 29446 ns
2022/03/08 11:55:31 stats.go:46: INFO : mean: 0 ms, median: 0 ms, max: 0 ms, min: 0 ms, p99.9: 0 ms
```

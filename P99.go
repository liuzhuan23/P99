package P99

import (
	"github.com/rpcxio/rpcx-benchmark/stat"
	"github.com/smallnest/rpcx/log"
	"go.uber.org/ratelimit"
	"sync"
	"sync/atomic"
	"time"
)

type P99CallBack func() int

type P99Stat struct {
	total int
	n     int
	m     int
	c     []P99IFace
}

type P99IFace struct {
	pInit P99CallBack
	pRun  P99CallBack
}

func NewP99Stat(t1 int, n1 int) *P99Stat {
	return &P99Stat{total: t1, n: n1, m: t1 / n1}
}

func (p *P99Stat) P99Register(init P99CallBack, run P99CallBack) {
	var pf P99IFace = P99IFace{pInit: init, pRun: run}
	p.c = append(p.c, pf)
}

func (p *P99Stat) P99IfCount() int {
	return len(p.c)
}

func (p *P99Stat) P99Run() {
	var rl ratelimit.Limiter

	// 总请求数
	var trans uint64

	// 返回正常的总请求数
	var transOK uint64

	// 等待所有事务测试完成的信号灯，信号数量是n*m
	var wg sync.WaitGroup
	wg.Add(p.n * p.m)

	// 每个goroutine的耗时记录
	d := make([][]int64, p.n, p.n)

	// 栅栏，控制客户端同时开始测试
	var startWg sync.WaitGroup
	// +1 是因为有一个goroutine用来记录开始时间
	startWg.Add(p.n + 1)

	startTime := time.Now().UnixNano()
	go func() {
		startWg.Done()
		startWg.Wait()
		startTime = time.Now().UnixNano()
	}()

	for i := 0; i < p.n; i++ {
		dt := make([]int64, 0, p.m)
		d = append(d, dt)

		// call 'init' handle hook
		if p.c[i].pInit != nil {
			p.c[i].pInit()
		}

		go func(i int) {
			code := -1
			for j := 0; j < p.m; j++ {
				if rl != nil {
					rl.Take()
				}

				t := time.Now().UnixNano()

				// call 'run' handle hook
				if p.c[i].pRun != nil {
					p.c[i].pRun()
				}

				// 等待时间+服务时间，等待时间是客户端调度的等待时间以及服务端读取请求、调度的时间，服务时间是请求被服务处理的实际时间
				t = time.Now().UnixNano() - t
				d[i] = append(d[i], t)
				if code != -1 {
					atomic.AddUint64(&transOK, 1)
				}
				atomic.AddUint64(&trans, 1)
				wg.Done()
			}
		}(i)
	}

	//等待信号灯
	wg.Wait()

	// 统计与输出
	log.Infof("P99客户端数量: %d\n", p.n)
	log.Infof("P99每个客户端的并发请求数: %d\n", p.m)
	stat.Stats(startTime, p.total, d, trans, transOK)
}

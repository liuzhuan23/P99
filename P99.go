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

// pcb钩子是必须注册的，分成2个情况，1 未注册，使用者自己在事务routine内实现复用资源 2 注册，P99在发起routine前帮你复用资源
// pin钩子是可以不进行注册的，如果未注册这个钩子，那就需要使用者在go func routine内自己复用套接字、句柄一类的资源变量
type P99Stat struct {
	pcb   P99CallBack
	pin   P99CallBack
	total int
	n     int
	m     int
}

func NewP99Stat(t1 int, n1 int) *P99Stat {
	return &P99Stat{total: t1, n: n1, m: t1 / n1, pcb: nil, pin: nil}
}

func (p *P99Stat) InitState(f P99CallBack) {
	p.pin = f
}

func (p *P99Stat) Register(f P99CallBack) {
	p.pcb = f
}

func (p *P99Stat) Run() {
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
		if p.pin != nil {
			if err := p.pin(); err == -1 {
				log.Infof("P99无法初始化pin钩子，错误代码: %d，进入资源句柄复用模式\n", err)
			}
		} else {
			log.Infof("P99未发现pin钩子，进入资源句柄复用模式\n")
		}
		go func(i int) {
			code := -1
			for j := 0; j < p.m; j++ {
				if rl != nil {
					rl.Take()
				}

				t := time.Now().UnixNano()

				// 调用钩子运行服务测试
				if p.pcb != nil {
					code = p.pcb()
				} else {
					log.Infof("P99无法调用pcb钩子，错误代码: %d\n", code)
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

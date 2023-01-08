package utils

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

// GeneratorFactor Generate 生成24位订单号
// 前面17位代表时间精确到毫秒，中间3位代表进程id，最后4位代表序号
func generatorFactor() func() string {
	var counter int64 = 0
	return func() string {
		t := time.Now()
		s := t.Format("20060102150405")
		m := t.UnixNano()/1e6 - t.UnixNano()/1e9*1e3
		ms := sup(m, 3)
		p := os.Getpid() % 1000
		ps := sup(int64(p), 3)
        i := atomic.AddInt64(&counter, 1)
		r := i % 10000
		rs := sup(r, 4)
		n := fmt.Sprintf("%s%s%s%s", ms, s, ps, rs)
		return n
	}
}

func sup(i int64, n int) string {
	m := fmt.Sprintf("%d", i)
	for len(m) < n {
		m = fmt.Sprintf("0%s", m)
	}
	return m
}


var (
    ShareBillIDGenerate = generatorFactor()
    OrderIDGenerate = generatorFactor()
)

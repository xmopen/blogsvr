package report

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/xmopen/blogsvr/internal/models/reportmod"
	"github.com/xmopen/golib/pkg/xlogging"

	"github.com/xmopen/golib/pkg/xgoroutine"
)

// TODO: 1、高并发情况下,内存定时刷新.
// TODO: 2、服务优雅推出.

// DBReport  reportmod to redis.
// 每分钟刷新上报内容到DB中.
type DBReport struct {
	reports *sync.Map
	lock    *sync.Mutex

	xlog *xlogging.Entry
}

// NewDBReport 初始化redis上报实例.
// 上层已经是单例模式.
func NewDBReport() IReport {
	dbReport := &DBReport{
		reports: &sync.Map{},
		lock:    &sync.Mutex{},
		xlog:    xlogging.Tag("report.db"),
	}
	dbReport.ticker()
	return dbReport
}

// ticker 定时刷新内存上报的数据到redis中.
func (r *DBReport) ticker() {
	xgoroutine.SafeGoroutine(func() {
		ticker := time.NewTicker(1 * time.Minute)
		for {
			select {
			case <-ticker.C:
				r.Flush()
			}
		}
	})
}

// Flush 刷新上报数据到DB.
func (r *DBReport) Flush() {
	r.reports.Range(func(key, value any) bool {
		incrKey, ok := key.(ReportType)
		if !ok {
			return true
		}
		incrValue, ok := value.(*atomic.Int64)
		if !ok {
			return true
		}
		if incrValue.Load() <= 0 {
			return true
		}
		if err := reportmod.Incr(int(incrKey), int(incrValue.Load())); err != nil {
			r.xlog.Errorf("reportmod incr err:[%+v] key:[%+v]", err, incrKey)
		} else {
			r.lock.Lock()
			r.reports.Store(incrKey, &atomic.Int64{})
			r.lock.Unlock()
		}
		return true
	})
}

// Report 递增上报.
func (r *DBReport) Report(key ReportType) {
	_, ok := r.reports.Load(key)
	if !ok {
		r.lock.Lock()
		r.reports.Store(key, &atomic.Int64{})
		r.lock.Unlock()
	}
	itr, _ := r.reports.Load(key)
	incr := itr.(*atomic.Int64)
	incr.Add(1)
}

// ReportWithValue 上报value.
func (r *DBReport) ReportWithValue(key ReportType, value any) {
	panic("no implement ReportWithValue")
}

// Package reportmod  上报.
package report

import (
	"sync"
)

// ReportType  上报类型.
type ReportType int

// 上报类型定义.
const (
	// ReportTypeIndex 首页上报.
	ReportTypeIndex ReportType = iota
	// ReportTypeArticleInfo 详情页上报.
	ReportTypeArticleInfo
)

var (
	reportInstance IReport
	reportOnce     sync.Once
)

// IReport  reportmod interface.
type IReport interface {
	// Report 递增上报.
	Report(key ReportType)
	// ReportWithValue reportmod with value.
	ReportWithValue(key ReportType, value any)
}

// Report  获取上报实例.
func Report() IReport {
	if reportInstance == nil {
		if reportInstance == nil {
			reportOnce.Do(func() {
				reportInstance = NewDBReport()
			})
		}
	}
	return reportInstance
}

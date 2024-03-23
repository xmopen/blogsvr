package reportmod

import "github.com/xmopen/blogsvr/internal/config"

const visitReportTableName = "t_report_visit"

// VisitReport 访问上报
type VisitReport struct {
	ID   int    `json:"id" gorm:"column:id"`
	Path string `json:"path" gorm:"column:path"`
	IP   string `json:"ip" gorm:"column:ip"`
}

// Visit 访问上报。
func Visit(visit *VisitReport) error {
	return config.BlogsDataBase().Table(visitReportTableName).Create(visit).Error
}

package reportmod

import (
	"fmt"

	"github.com/xmopen/blogsvr/internal/config"
)

const reportTableName = "t_report"

var incrSQL = `update t_report set incr_value = incr_value + %d where id = %d limit 1;`

// Report  report struct.
type Report struct {
	ID        int    `json:"id" gorm:"column:id"`
	Desc      string `json:"desc" gorm:"column:desc"`
	IncrValue int    `json:"incr" gorm:"column:incr_value"`
	Value     string `json:"value" gorm:"column:value"`
}

// Incr incr.
func Incr(id, incr int) error {
	return config.BlogsDataBase().Exec(fmt.Sprintf(incrSQL, incr, id)).Error
}

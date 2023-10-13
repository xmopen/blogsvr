package archivemod

import (
	"errors"

	"github.com/xmopen/blogsvr/internal/config"
	"github.com/xmopen/commonlib/pkg/database/xmarchive"
	"gorm.io/gorm"
)

const tXMBlogsArchiveTableName = "t_xm_blogs_archive"

// GetArchiveList return archive list
func GetArchiveList() ([]*xmarchive.XMBlogsArchive, error) {
	resultList := make([]*xmarchive.XMBlogsArchive, 0)
	err := config.BlogsDataBase().Table(tXMBlogsArchiveTableName).Where("status = 1").Find(&resultList).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return resultList, nil
}

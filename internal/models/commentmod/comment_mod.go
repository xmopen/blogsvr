package commentmod

import (
	"errors"

	"github.com/xmopen/blogsvr/internal/config"
	"github.com/xmopen/commonlib/pkg/database/xmcomment"
	"gorm.io/gorm"
)

const tXMCommentTableName = "t_xm_comment"

// GetCommentListByArticleID get comment list by article id
func GetCommentListByArticleID(articleID int) ([]*xmcomment.XMComment, error) {
	resultList := make([]*xmcomment.XMComment, 0)
	err := config.BlogsDataBase().Table(tXMCommentTableName).Where("article_id=?", articleID).Find(&resultList).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return resultList, nil
}

// CreateXMComment create comment
func CreateXMComment(comment *xmcomment.XMComment) error {
	return config.BlogsDataBase().Table(tXMCommentTableName).Create(comment).Error
}

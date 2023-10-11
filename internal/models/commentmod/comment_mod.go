package commentmod

import (
	"errors"

	"github.com/xmopen/blogsvr/internal/config"
	"github.com/xmopen/commonlib/pkg/database/xmcomment"
	"gorm.io/gorm"
)

// ArticleCommentStatus 评论状态
type ArticleCommentStatus int

const (
	// ArticleCommentStatusOfDown 下架评论
	ArticleCommentStatusOfDown ArticleCommentStatus = iota
	// ArticleCommentStatusOfUp 上架评论
	ArticleCommentStatusOfUp
)

const tXMCommentTableName = "t_xm_comment"

// Comment comment struct
type Comment struct {
	Username    string `json:"username"`
	Icon        string `json:"icon"`
	CommentTime string `json:"comment_time"`
	*xmcomment.XMComment
}

// GetCommentListByArticleID get comment list by article id
func GetCommentListByArticleID(articleID int) ([]*xmcomment.XMComment, error) {
	resultList := make([]*xmcomment.XMComment, 0)
	err := config.BlogsDataBase().Table(tXMCommentTableName).Where("article_id=? and `status`=?", articleID,
		ArticleCommentStatusOfUp).Find(&resultList).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return resultList, nil
}

// CreateXMComment create a comment for article
func CreateXMComment(comment *xmcomment.XMComment) error {
	return config.BlogsDataBase().Table(tXMCommentTableName).Create(comment).Error
}

package articlemod

import (
	"errors"

	"github.com/xmopen/blogsvr/internal/config"
	"gorm.io/gorm"
)

const articleTableName = "blogs_article"

// ArticlePublishType 文章发布状态.
type ArticlePublishType int8

// 文章发布状态.
const (
	// ArticlePublishNotPublish 文章未发布.
	ArticlePublishNotPublish ArticlePublishType = iota
	// ArticlePublishOldPublished V1文章发布状态.
	ArticlePublishOldPublished
	// ArticlePublish 文章发布.
	ArticlePublish
)

// Article 文章.
type Article struct {
	ID      int    `json:"id" gorm:"column:id"`
	Title   string `json:"title" gorm:"column:title"`
	Time    string `json:"time" gorm:"column:time"`
	Author  string `json:"author" gorm:"column:author"`
	Content string `json:"content" gorm:"column:content"`
	SubHead string `json:"sub_head" gorm:"column:subhead"`
	Img     string `json:"img" gorm:"column:img"`
	TypeID  string `json:"type_id" gorm:"column:type_id"` // 分类ID.
	Publish int    `json:"publish" gorm:"column:publish"` // 是否发布.
}

// AllArticles 获取所有文章信息.
func AllArticles() ([]*Article, error) {
	articleList := make([]*Article, 0)
	result := config.BlogsDataBase().Table(articleTableName).Where("publish = ?", ArticlePublish).Find(&articleList)
	if result.Error != nil {
		// ErrRecordNotFound
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return articleList, nil
		}
	}
	return articleList, nil
}

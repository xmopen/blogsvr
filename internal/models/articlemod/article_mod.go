package articlemod

import (
	"errors"
	"time"

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
	// ArticlePublished published.
	ArticlePublished
)

// Article 文章.
type Article struct {
	ID          int       `json:"id" gorm:"column:id"`
	TypeID      int       `json:"type_id" gorm:"column:type_id"` // TypeID 分类ID.
	Publish     int       `json:"publish" gorm:"column:publish"` // Publish 是否发布.
	ReadCount   int       `json:"read_count" gorm:"column:read_count"`
	Title       string    `json:"title" gorm:"column:title"`
	Time        string    `json:"time" gorm:"column:time"`
	Author      string    `json:"author" gorm:"column:author"`
	Content     string    `json:"content" gorm:"column:content"`
	SubHead     string    `json:"sub_head" gorm:"column:subhead"`
	Img         string    `json:"img" gorm:"column:img"`
	PublishTime time.Time `json:"publish_time" gorm:"column:publish_time"`
	Type        string    `json:"type" gorm:"-"` // Type 分类
}

// AllArticles 获取所有文章信息.
func AllArticles() ([]*Article, error) {
	articleList := make([]*Article, 0)
	result := config.BlogsDataBase().Table(articleTableName).Where("publish = ?", ArticlePublished).
		Order("id desc").Find(&articleList)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return articleList, nil
		}
		return nil, result.Error
	}
	for _, article := range articleList {
		if article.Time != "" {
			continue
		}
		article.Time = article.PublishTime.Format(time.DateTime)
	}
	return articleList, nil
}

// UpdateArticleReadCount 更新文章ReadCount
func UpdateArticleReadCount(articleID int) error {
	return config.BlogsDataBase().Table(articleTableName).Where("id = ?", articleID).
		Update("read_count", gorm.Expr("read_count + ?", 1)).Error
}

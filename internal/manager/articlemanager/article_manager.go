// Package articlemanager 文章管理器.
package articlemanager

import (
	"sync"
)

var (
	articleManagerInstance     *ArticleManager
	articleManagerInstanceOnce sync.Once
)

// ArticleManager 文章管理器.
type ArticleManager struct {
}

// Manager 返回文章管理器. articlemanager.Manager()
func Manager() *ArticleManager {
	articleManagerInstanceOnce.Do(func() {
		articleManagerInstance = &ArticleManager{}
	})
	return articleManagerInstance
}

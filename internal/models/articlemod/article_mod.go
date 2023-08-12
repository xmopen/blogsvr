package articlemod

// Article 文章.
type Article struct {
	ID      int    `json:"id" gorm:"column:id"`
	Title   string `json:"title" gorm:"column:title"`
	Time    string `json:"time" gorm:"column:time"`
	Author  string `json:"author" gorm:"column:author"`
	Content string `json:"content" gorm:"column:content"`
	SubHead string `json:"sub_head" gorm:"column:subhead"`
	Type    string `json:"type" gorm:"column:type"`
	Img     string `json:"img" gorm:"column:img"`
}

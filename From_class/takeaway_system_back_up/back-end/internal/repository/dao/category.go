package dao

type Category struct {
	CategoryID  int64  `gorm:"primary_key;auto_increment;comment:分类唯一标识"`
	Name        string `gorm:"type:varchar(100);not null;comment:分类名称"`
	Description string `gorm:"type:text;comment:分类描述"`
}

package model

// User model
type User struct {
	ID       int    `gorm:"column:user_id;type:int;primaryKey;autoIncrement"`
	Username string `gorm:"column:username;type:varchar(50);notnull;uniqueIndex"`
	Name     string `gorm:"column:name;type:varchar(50)"`
	Role     string `gorm:"column:role;type:varchar(20);notnull"`
	Password string `gorm:"column:password;type:varchar(255)"`
}

package models

import (
	"time"

	"gorm.io/gorm"
)

// UserRole represents the role of a user
type UserRole string

const (
	RoleAdmin UserRole = "admin" //管理员
	RoleUser  UserRole = "user"  //普通用户
)

// UserStatus represents the status of a user
type UserStatus string

const (
	StatusActive UserStatus = "active" //正常
	StatusBanned UserStatus = "banned" //封禁
)

type User struct {
	gorm.Model
	Username  string     `gorm:"unique;not null"`
	Password  string     `gorm:"not null"`
	Nickname  string     `gorm:"size:50"`
	AvatarURL string     `gorm:"size:255"`
	Role      UserRole   `gorm:"type:varchar(10);not null;default:'user'"`   //用户角色
	Status    UserStatus `gorm:"type:varchar(10);not null;default:'active'"` //用户状态 是否封禁
	Guides    []TravelGuide
}

type TravelGuide struct {
	gorm.Model
	Title       string    `gorm:"not null;type:varchar(255)"`
	Content     string    `gorm:"type:text;not null"`
	Images      string    `gorm:"type:text"` // Store image URLs as JSON string
	UserID      uint      `gorm:"not null"`
	User        User      `gorm:"foreignKey:UserID"`
	PublishedAt time.Time `gorm:"not null"`
	Tags        []Tag     `gorm:"many2many:guide_tags;"`
}

type Tag struct {
	gorm.Model
	Name   string        `gorm:"unique;not null"`
	Guides []TravelGuide `gorm:"many2many:guide_tags;"`
}

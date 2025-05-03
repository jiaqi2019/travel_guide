package models

import (
	"time"
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
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Username  string    `gorm:"unique;not null;size:50"`
	Password  string    `gorm:"not null;size:255"`
	Nickname  string    `gorm:"not null;size:100"`
	AvatarURL string    `gorm:"size:255"`
	Role      UserRole  `gorm:"type:enum('admin','user');not null;default:'user'"`
	Status    UserStatus `gorm:"type:enum('active','banned');not null;default:'active'"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt *time.Time
	Guides    []TravelGuide
	Tags      []Tag     `gorm:"many2many:user_tags;joinForeignKey:user_id;joinReferences:tag_id"`
}

type TravelGuide struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Title       string    `gorm:"not null;size:255"`
	Content     string    `gorm:"not null;type:text"`
	Images      string    `gorm:"type:text"`
	UserID      uint      `gorm:"not null"`
	User        User      `gorm:"foreignKey:UserID"`
	PublishedAt time.Time `gorm:"not null"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt   *time.Time
	Tags        []Tag     `gorm:"many2many:guide_tags;joinForeignKey:guide_id;joinReferences:tag_id"`
}

type Tag struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"unique;not null;size:50"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt *time.Time
	Guides    []TravelGuide `gorm:"many2many:guide_tags;joinForeignKey:tag_id;joinReferences:guide_id"`
}

package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.DBConfig.User,
		AppConfig.DBConfig.Password,
		AppConfig.DBConfig.Host,
		AppConfig.DBConfig.Port,
		AppConfig.DBConfig.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// 创建数据库（如果不存在）
	err = db.Exec("CREATE DATABASE IF NOT EXISTS travel_guide DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci").Error
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %v", err)
	}

	// 使用数据库
	err = db.Exec("USE travel_guide").Error
	if err != nil {
		return nil, fmt.Errorf("failed to use database: %v", err)
	}

	// 创建表（如果不存在）
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(50) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			nickname VARCHAR(100) NOT NULL,
			avatar_url VARCHAR(255),
			role ENUM('admin', 'user') NOT NULL DEFAULT 'user',
			status ENUM('active', 'banned') NOT NULL DEFAULT 'active',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL,
			INDEX idx_username (username),
			INDEX idx_nick_name (nickname)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create users table: %v", err)
	}

	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tags (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(50) NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL,
			INDEX idx_name (name)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create tags table: %v", err)
	}

	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS travel_guide (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			images TEXT,
			user_id BIGINT UNSIGNED NOT NULL,
			published_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FULLTEXT INDEX idx_ft_title_content (title, content),
			INDEX idx_user_id (user_id),
			INDEX idx_published_at (published_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create travel_guide table: %v", err)
	}

	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS guide_tags (
			guide_id BIGINT UNSIGNED NOT NULL,
			tag_id BIGINT UNSIGNED NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (guide_id, tag_id),
			FOREIGN KEY (guide_id) REFERENCES travel_guide(id),
			FOREIGN KEY (tag_id) REFERENCES tags(id),
			INDEX idx_tag_id (tag_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create guide_tags table: %v", err)
	}

	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS user_tags (
			user_id BIGINT UNSIGNED NOT NULL,
			tag_id BIGINT UNSIGNED NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, tag_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (tag_id) REFERENCES tags(id),
			INDEX idx_tag_id (tag_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create user_tags table: %v", err)
	}

	// 插入初始数据（如果不存在）
	err = db.Exec(`
		INSERT IGNORE INTO tags (name) VALUES 
		('自然风光'),
		('城市探索'),
		('海岛度假'),
		('乡村田园'),
		('高原雪山'),
		('沙漠戈壁'),
		('森林徒步'),
		('草原牧场');
	`).Error
	if err != nil {
		return nil, fmt.Errorf("failed to insert initial tags: %v", err)
	}

	err = db.Exec(`
		INSERT IGNORE INTO users (username, password, nickname, avatar_url, role) VALUES 
		('admin', '$2a$10$CPmq3FF9E9Vf622Tt5bSmOAKLATwZraADO3haRCcMGcHeU7lfJakC', 'admin', 'https://example.com/avatar1.jpg', 'admin'),
		('travel_lover', '$2a$10$CPmq3FF9E9Vf622Tt5bSmOAKLATwZraADO3haRCcMGcHeU7lfJakC', '旅行爱好者', 'https://example.com/avatar1.jpg', 'user'),
		('food_explorer', '$2a$10$CPmq3FF9E9Vf622Tt5bSmOAKLATwZraADO3haRCcMGcHeU7lfJakC', '美食探索家', 'https://example.com/avatar2.jpg', 'user'),
		('nature_seeker', '$2a$10$CPmq3FF9E9Vf622Tt5bSmOAKLATwZraADO3haRCcMGcHeU7lfJakC', '自然追寻者', 'https://example.com/avatar3.jpg', 'user'),
		('city_wanderer', '$2a$10$CPmq3FF9E9Vf622Tt5bSmOAKLATwZraADO3haRCcMGcHeU7lfJakC', '城市漫游者', 'https://example.com/avatar4.jpg', 'user'),
		('adventure_seeker', '$2a$10$CPmq3FF9E9Vf622Tt5bSmOAKLATwZraADO3haRCcMGcHeU7lfJakC', '冒险家', 'https://example.com/avatar5.jpg', 'user');
	`).Error
	if err != nil {
		return nil, fmt.Errorf("failed to insert initial users: %v", err)
	}

	err = db.Exec(`
		INSERT IGNORE INTO user_tags (user_id, tag_id) 
		SELECT u.id, t.id 
		FROM users u 
		CROSS JOIN tags t 
		WHERE 
			(u.username = 'travel_lover' AND t.name IN ('自然风光', '城市探索', '海岛度假')) OR
			(u.username = 'food_explorer' AND t.name IN ('城市探索', '乡村田园')) OR
			(u.username = 'nature_seeker' AND t.name IN ('自然风光', '高原雪山', '森林徒步', '草原牧场')) OR
			(u.username = 'city_wanderer' AND t.name IN ('城市探索', '海岛度假')) OR
			(u.username = 'adventure_seeker' AND t.name IN ('高原雪山', '沙漠戈壁', '森林徒步'));
	`).Error
	if err != nil {
		return nil, fmt.Errorf("failed to insert initial user_tags: %v", err)
	}

	return db, nil
} 
-- 创建数据库
CREATE DATABASE IF NOT EXISTS travel_guide DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE travel_guide;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    nickname VARCHAR(100) NOT NULL,
    avatar_url VARCHAR(255),
    role ENUM('admin', 'user') NOT NULL DEFAULT 'user' COMMENT '用户角色：admin-管理员，user-普通用户',
    status ENUM('active', 'banned') NOT NULL DEFAULT 'active' COMMENT '用户状态：active-正常，banned-封禁',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_username (username),
    INDEX idx_nick_name (nickname)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 标签表
CREATE TABLE IF NOT EXISTS tags (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 旅游攻略表
CREATE TABLE IF NOT EXISTS travel_guides (
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

-- 攻略-标签关联表
CREATE TABLE IF NOT EXISTS guide_tags (
    guide_id BIGINT UNSIGNED NOT NULL,
    tag_id BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (guide_id, tag_id),
    FOREIGN KEY (guide_id) REFERENCES travel_guides(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id),
    INDEX idx_tag_id (tag_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 用户-标签关联表
CREATE TABLE IF NOT EXISTS user_tags (
    user_id BIGINT UNSIGNED NOT NULL,
    tag_id BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, tag_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id),
    INDEX idx_tag_id (tag_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 插入初始标签
INSERT IGNORE INTO tags (name) VALUES 
-- 旅行类型
('自然风光'),
('城市探索'),
('海岛度假'),
('乡村田园'),
('高原雪山'),
('沙漠戈壁'),
('森林徒步'),
('草原牧场');

-- 插入模拟用户数据
INSERT IGNORE INTO users (username, password, nickname, avatar_url, role) VALUES 
('admin', '$2a$10$CPmq3FF9E9Vf622Tt5bSmOAKLATwZraADO3haRCcMGcHeU7lfJakC', 'admin', 'https://api.dicebear.com/7.x/initials/svg?seed=admin', 'admin'),
('travel_lover', '$2a$10$CPmq3FF9E9Vf622Tt5bSmOAKLATwZraADO3haRCcMGcHeU7lfJakC', '旅行爱好者', 'https://api.dicebear.com/7.x/initials/svg?seed=旅行爱好者', 'user'),
('food_explorer', '$2a$10$CPmq3FF9E9Vf622Tt5bSmOAKLATwZraADO3haRCcMGcHeU7lfJakC', '美食探索家', 'https://api.dicebear.com/7.x/initials/svg?seed=美食探索家', 'user'),
('nature_seeker', '$2a$10$CPmq3FF9E9Vf622Tt5bSmOAKLATwZraADO3haRCcMGcHeU7lfJakC', '自然追寻者', 'https://api.dicebear.com/7.x/initials/svg?seed=自然追寻者', 'user'),
('city_wanderer', '$2a$10$CPmq3FF9E9Vf622Tt5bSmOAKLATwZraADO3haRCcMGcHeU7lfJakC', '城市漫游者', 'https://api.dicebear.com/7.x/initials/svg?seed=城市漫游者', 'user'),
('adventure_seeker', '$2a$10$CPmq3FF9E9Vf622Tt5bSmOAKLATwZraADO3haRCcMGcHeU7lfJakC', '冒险家', 'https://api.dicebear.com/7.x/initials/svg?seed=冒险家', 'user');

-- 插入用户标签关联数据
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

-- 插入攻略标签关联数据
INSERT IGNORE INTO guide_tags (guide_id, tag_id)
SELECT g.id, t.id
FROM travel_guides g
CROSS JOIN tags t
WHERE 
    (g.title = '探索巴厘岛的神秘之美' AND t.name IN ('海岛度假', '自然风光')) OR
    (g.title = '成都美食之旅：舌尖上的川菜盛宴' AND t.name IN ('城市探索', '乡村田园')) OR
    (g.title = '西藏高原之旅：探索世界屋脊' AND t.name IN ('高原雪山', '自然风光')) OR
    (g.title = '东京城市探索：传统与现代的完美融合' AND t.name IN ('城市探索')) OR
    (g.title = '撒哈拉沙漠探险：感受大自然的壮丽' AND t.name IN ('沙漠戈壁', '自然风光'));

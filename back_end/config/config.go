package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConfig     DBConfig
	OSSConfig    OSSConfig
	ServerConfig ServerConfig
	JWTConfig    JWTConfig
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type OSSConfig struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
}

type ServerConfig struct {
	Port int
}

type JWTConfig struct {
	SecretKey string
	ExpiresIn int64
}

var AppConfig Config

func LoadConfig() error {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		// 如果找不到 .env 文件，继续使用环境变量
	}

	// 数据库配置
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))
	AppConfig.DBConfig = DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     dbPort,
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "travel_guide"),
	}

	// OSS配置
	AppConfig.OSSConfig = OSSConfig{
		Endpoint:        getEnv("OSS_ENDPOINT", ""),
		AccessKeyID:     getEnv("OSS_ACCESS_KEY_ID", ""),
		AccessKeySecret: getEnv("OSS_ACCESS_KEY_SECRET", ""),
		BucketName:      getEnv("OSS_BUCKET_NAME", ""),
	}

	// 服务器配置
	serverPort, _ := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	AppConfig.ServerConfig = ServerConfig{
		Port: serverPort,
	}

	// JWT配置
	expiresIn, _ := strconv.ParseInt(getEnv("JWT_EXPIRES_IN", "86400"), 10, 64)
	AppConfig.JWTConfig = JWTConfig{
		SecretKey: getEnv("JWT_SECRET_KEY", ""),
		ExpiresIn: expiresIn,
	}

	return nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
} 
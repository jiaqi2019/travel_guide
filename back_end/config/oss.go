package config

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var OSSClient *oss.Client

func InitOSS() error {
	// 创建OSSClient实例
	client, err := oss.New(
		AppConfig.OSSConfig.Endpoint,
		AppConfig.OSSConfig.AccessKeyID,
		AppConfig.OSSConfig.AccessKeySecret,
	)
	if err != nil {
		return err
	}

	OSSClient = client
	return nil
} 
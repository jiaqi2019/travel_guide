package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"travel_guide/config"
	"travel_guide/types"
	"travel_guide/utils/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadController struct{}

func NewUploadController() *UploadController {
	return &UploadController{}
}

// UploadImage 处理图片上传
func (uc *UploadController) UploadImage(c *gin.Context) {
	logger.InfoLogger.Println("开始处理图片上传")

	// 获取上传的文件
	file, err := c.FormFile("image")
	if err != nil {
		logger.ErrorLogger.Printf("获取上传文件失败: %v", err)
		c.JSON(http.StatusOK, types.ErrorResponse(1, "请选择要上传的图片"))
		return
	}
	logger.InfoLogger.Printf("接收到文件: %s, 大小: %d bytes", file.Filename, file.Size)

	// 检查文件类型
	ext := filepath.Ext(file.Filename)
	if !isValidImageType(ext) {
		logger.ErrorLogger.Printf("不支持的文件类型: %s", ext)
		c.JSON(http.StatusOK, types.ErrorResponse(1, "不支持的文件类型"))
		return
	}

	// 生成唯一的文件名
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	logger.InfoLogger.Printf("生成的文件名: %s", fileName)

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "upload-*"+ext)
	if err != nil {
		logger.ErrorLogger.Printf("创建临时文件失败: %v", err)
		c.JSON(http.StatusOK, types.ErrorResponse(1, "创建临时文件失败"))
		return
	}
	defer os.Remove(tempFile.Name())
	logger.InfoLogger.Printf("创建临时文件: %s", tempFile.Name())

	// 保存文件到临时目录
	if err := c.SaveUploadedFile(file, tempFile.Name()); err != nil {
		logger.ErrorLogger.Printf("保存上传文件失败: %v", err)
		c.JSON(http.StatusOK, types.ErrorResponse(1, "保存文件失败"))
		return
	}

	// 获取OSS bucket
	bucketName := os.Getenv("OSS_BUCKET_NAME")
	logger.InfoLogger.Printf("准备上传到OSS，Bucket: %s", bucketName)

	bucket, err := config.OSSClient.Bucket(bucketName)
	if err != nil {
		logger.ErrorLogger.Printf("获取OSS bucket失败: %v", err)
		c.JSON(http.StatusOK, types.ErrorResponse(1, "获取OSS bucket失败"))
		return
	}

	// 上传到OSS
	objectName := fmt.Sprintf("images/%s/%s", time.Now().Format("2006/01/02"), fileName)
	logger.InfoLogger.Printf("开始上传文件到OSS，对象名: %s", objectName)

	err = bucket.PutObjectFromFile(objectName, tempFile.Name())
	if err != nil {
		logger.ErrorLogger.Printf("上传到OSS失败: %v", err)
		c.JSON(http.StatusOK, types.ErrorResponse(1, "上传图片失败"))
		return
	}

	// 生成访问URL
	url := fmt.Sprintf("https://%s.%s/%s", bucketName, config.OSSClient.Config.Endpoint, objectName)
	logger.InfoLogger.Printf("文件上传成功，访问URL: %s", url)

	c.JSON(http.StatusOK, types.SuccessResponse(
		gin.H{"url": url},
		"上传成功",
	))
}

// isValidImageType 检查文件类型是否为图片
func isValidImageType(ext string) bool {
	validTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	return validTypes[ext]
}

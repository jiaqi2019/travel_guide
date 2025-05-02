package controllers

import (
	"net/http"

	"travel_guide/models"
	"travel_guide/types"
	"travel_guide/utils/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TagController struct {
	db *gorm.DB
}

type TagResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type TagListResponse struct {
	List    []types.TagResponse `json:"list"`
	HasMore bool                `json:"has_more"`
}

func NewTagController(db *gorm.DB) *TagController {
	return &TagController{db: db}
}

func (tc *TagController) GetAllTags(c *gin.Context) {
	var tags []models.Tag
	if err := tc.db.Find(&tags).Error; err != nil {
		logger.ErrorLogger.Printf("获取标签失败: %v", err)
		c.JSON(http.StatusOK, types.ErrorResponse(1, "获取标签失败"))
		return
	}

	var response []types.TagResponse
	for _, tag := range tags {
		response = append(response, types.TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	c.JSON(http.StatusOK, types.SuccessResponse(response, "获取标签成功"))
}

// GetRelatedTags 获取搜索词相关的标签
func (tc *TagController) GetRelatedTags(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusOK, types.SuccessResponse(
			TagListResponse{
				List:    []types.TagResponse{},
				HasMore: false,
			},
			"获取相关标签成功",
		))
		return
	}

	// 获取相关标签
	var relatedTags []models.Tag
	var totalCount int64

	// 先获取总数
	tc.db.Model(&models.Tag{}).
		Joins("JOIN guide_tags ON guide_tags.tag_id = tags.id").
		Joins("JOIN travel_guides ON travel_guides.id = guide_tags.travel_guide_id").
		Where("travel_guides.title LIKE ? OR travel_guides.content LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Group("tags.id").
		Count(&totalCount)

	// 获取当前页数据
	tc.db.Model(&models.Tag{}).
		Joins("JOIN guide_tags ON guide_tags.tag_id = tags.id").
		Joins("JOIN travel_guides ON travel_guides.id = guide_tags.travel_guide_id").
		Where("travel_guides.title LIKE ? OR travel_guides.content LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Group("tags.id").
		Order("COUNT(*) DESC").
		Limit(10).
		Find(&relatedTags)

	// 转换标签响应格式
	tagResponses := make([]types.TagResponse, 0, len(relatedTags))
	for _, tag := range relatedTags {
		tagResponses = append(tagResponses, types.TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	c.JSON(http.StatusOK, types.SuccessResponse(
		TagListResponse{
			List:    tagResponses,
			HasMore: totalCount > 10,
		},
		"获取相关标签成功",
	))
}

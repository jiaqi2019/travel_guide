package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"travel_guide/models"
	"travel_guide/types"
	"travel_guide/utils/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 转换单个攻略
func toGuideResponse(guide models.TravelGuide) types.GuideResponse {
	// 解析 images
	var images []string
	_ = json.Unmarshal([]byte(guide.Images), &images)

	// 转换 tags
	var tags []types.TagResponse
	for _, tag := range guide.Tags {
		tags = append(tags, types.TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	// 转换用户信息
	userResponse := types.UserResponse{
		ID:        guide.User.ID,
		Username:  guide.User.Username,
		Nickname:  guide.User.Nickname,
		AvatarURL: guide.User.AvatarURL,
	}

	return types.GuideResponse{
		ID:          guide.ID,
		Title:       guide.Title,
		Content:     guide.Content,
		Images:      images,
		UserID:      guide.UserID,
		User:        userResponse,
		PublishedAt: guide.PublishedAt.Unix(),
		Tags:        tags,
	}
}

// 转换攻略列表
func toGuideResponseList(guides []models.TravelGuide) []types.GuideResponse {
	var responses []types.GuideResponse
	for _, guide := range guides {
		responses = append(responses, toGuideResponse(guide))
	}
	return responses
}

// 添加一个转换函数
func toCreateGuideResponse(guide models.TravelGuide) types.CreateGuideResponse {
	// 解析 images
	var images []string
	_ = json.Unmarshal([]byte(guide.Images), &images)

	// 转换 tags
	var tags []types.TagResponse
	for _, tag := range guide.Tags {
		tags = append(tags, types.TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	return types.CreateGuideResponse{
		ID:          guide.ID,
		Title:       guide.Title,
		Content:     guide.Content,
		Images:      images,
		PublishedAt: guide.PublishedAt.Unix(),
		Tags:        tags,
	}
}

type GuideController struct {
	db *gorm.DB
}

func NewGuideController(db *gorm.DB) *GuideController {
	return &GuideController{db: db}
}

type CreateGuideRequest struct {
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Images  []string `json:"images"`
	Tags    []string `json:"tags"`
}

// 创建攻略
func (gc *GuideController) CreateGuide(c *gin.Context) {
	logger.InfoLogger.Printf("开始创建新攻略")

	var req CreateGuideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.ErrorLogger.Printf("请求数据绑定失败: %v", err)
		c.JSON(http.StatusOK, types.ErrorResponse(1, "请求参数错误"))
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		logger.ErrorLogger.Printf("未找到用户ID，认证失败")
		c.JSON(http.StatusOK, types.ErrorResponse(1, "未授权"))
		return
	}

	// 检查用户状态
	var user models.User
	if err := gc.db.First(&user, userID).Error; err != nil {
		logger.ErrorLogger.Printf("用户不存在: %v", err)
		c.JSON(http.StatusOK, types.ErrorResponse(1, "用户不存在"))
		return
	}

	if user.Status == models.StatusBanned {
		logger.ErrorLogger.Printf("用户已被禁用，无法创建攻略")
		c.JSON(http.StatusOK, types.ErrorResponse(1, "用户已被禁用"))
		return
	}

	logger.InfoLogger.Printf("用户 %v 开始创建攻略", userID)

	imagesJSON, err := json.Marshal(req.Images)
	if err != nil {
		logger.ErrorLogger.Printf("图片JSON序列化失败: %v", err)
		c.JSON(http.StatusOK, types.ErrorResponse(1, "处理图片失败"))
		return
	}
	logger.InfoLogger.Printf("处理图片列表: %s", string(imagesJSON))

	guide := models.TravelGuide{
		Title:       req.Title,
		Content:     req.Content,
		Images:      string(imagesJSON),
		UserID:      userID.(uint),
		PublishedAt: time.Now(),
	}

	// Create or get tags
	var tags []models.Tag
	for _, tagName := range req.Tags {
		var tag models.Tag
		gc.db.FirstOrCreate(&tag, models.Tag{Name: tagName})
		tags = append(tags, tag)
	}

	guide.Tags = tags

	// 使用事务来确保数据一致性
	err = gc.db.Transaction(func(tx *gorm.DB) error {
		// 创建攻略
		if err := tx.Create(&guide).Error; err != nil {
			return err
		}

		// 查询完整信息，只加载标签信息
		if err := tx.Preload("Tags").First(&guide, guide.ID).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		logger.ErrorLogger.Printf("保存攻略失败: %v", err)
		c.JSON(http.StatusOK, types.ErrorResponse(1, "创建攻略失败"))
		return
	}

	logger.InfoLogger.Printf("攻略创建成功，ID: %v", guide.ID)
	c.JSON(http.StatusOK, types.SuccessResponse(
		toCreateGuideResponse(guide),
		"创建攻略成功",
	))
}

// 关键词查找攻略
func (gc *GuideController) SearchGuides(c *gin.Context) {
	keyword := c.Query("keyword")
	tag := c.Query("tag")
	offset, _ := c.GetQuery("offset")
	limit, _ := c.GetQuery("limit")

	// 设置默认值
	offsetInt := 0
	limitInt := 10
	if offset != "" {
		offsetInt, _ = strconv.Atoi(offset)
	}
	if limit != "" {
		limitInt, _ = strconv.Atoi(limit)
	}
	if limitInt <= 0 {
		limitInt = 10
	}

	logger.InfoLogger.Printf("搜索攻略 - 关键词: %s, 标签: %s, 偏移: %d, 限制: %d", keyword, tag, offsetInt, limitInt)

	// 构建基础查询
	query := gc.db.Model(&models.TravelGuide{}).
		Preload("User").
		Preload("Tags")

	// 添加关键词搜索条件
	if keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 如果有标签过滤条件，添加标签过滤
	if tag != "" {
		query = query.Joins("JOIN guide_tags ON guide_tags.guide_id = travel_guides.id").
			Joins("JOIN tags ON tags.id = guide_tags.tag_id").
			Where("tags.name = ?", tag)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 查询数据
	var guides []models.TravelGuide
	query.Offset(offsetInt).Limit(limitInt + 1).Find(&guides) // 多查询一条用于判断是否还有更多

	// 处理分页结果
	hasMore := false
	if len(guides) > limitInt {
		hasMore = true
		guides = guides[:limitInt] // 去掉多查询的一条
	}

	// 转换响应格式
	guideResponses := make([]types.GuideResponse, 0, len(guides))
	for _, guide := range guides {
		guideResponses = append(guideResponses, toGuideResponse(guide))
	}

	c.JSON(http.StatusOK, types.SuccessResponse(
		gin.H{
			"list":    guideResponses,
			"hasMore": hasMore,
			"total":   total,
		},
		"搜索成功",
	))
}

// GetRelatedTags 获取搜索词相关的标签
func (gc *GuideController) GetRelatedTags(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusOK, types.SuccessResponse(
			[]types.TagResponse{},
			"获取相关标签成功",
		))
		return
	}

	// 获取相关标签
	var relatedTags []models.Tag
	gc.db.Model(&models.Tag{}).
		Joins("JOIN guide_tags ON guide_tags.tag_id = tags.id").
		Joins("JOIN travel_guides ON travel_guides.id = guide_tags.travel_guide_id").
		Where("MATCH(travel_guides.title, travel_guides.content) AGAINST(? IN BOOLEAN MODE)", keyword).
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
		tagResponses,
		"获取相关标签成功",
	))
}

// 修改查询参数结构体
type GetGuidesRequest struct {
	Tag    string `form:"tag"`
	Offset int    `form:"offset,default=0"` // 改用 offset
	Limit  int    `form:"limit,default=10"` // page_size 改名为 limit
}

// 重命名为 GetGuides，因为是获取攻略列表
func (gc *GuideController) GetGuides(c *gin.Context) {
	var req GetGuidesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    400,
			Message: "参数错误",
		})
		return
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10 // 默认每页10条
	}

	query := gc.db.Model(&models.TravelGuide{}).
		Preload("User"). // 加载用户信息
		Preload("Tags")  // 加载标签信息

	// 如果有tag参数，添加tag过滤
	if req.Tag != "" {
		query = query.Joins("JOIN guide_tags ON guide_tags.guide_id = travel_guides.id").
			Joins("JOIN tags ON tags.id = guide_tags.tag_id").
			Where("tags.name = ?", req.Tag).
			Group("travel_guides.id") // 添加分组避免重复
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 查询数据
	var guides []models.TravelGuide
	query.Offset(req.Offset).Limit(limit + 1).Find(&guides) // 多查询一条用于判断是否还有更多

	hasMore := false
	if len(guides) > limit {
		hasMore = true
		guides = guides[:limit] // 去掉多查询的一条
	}

	// 转换响应格式
	guideResponses := make([]types.GuideResponse, 0, len(guides))
	for _, guide := range guides {
		guideResponses = append(guideResponses, toGuideResponse(guide))
	}

	c.JSON(http.StatusOK, types.Response{
		Code: 0,
		Data: types.PaginationResponse{
			List:    guideResponses,
			Total:   total,
			HasMore: hasMore,
		},
		Message: "success",
	})
}

// 获取单个攻略详情（重命名为 GetGuideDetail）
func (gc *GuideController) GetGuideDetail(c *gin.Context) {
	id := c.Param("id")
	logger.InfoLogger.Printf("获取攻略详情，ID: %s", id)

	var guide models.TravelGuide
	if err := gc.db.Preload("User").Preload("Tags").First(&guide, id).Error; err != nil {
		logger.ErrorLogger.Printf("获取攻略详情失败，ID %s: %v", id, err)
		c.JSON(http.StatusOK, types.ErrorResponse(1, "攻略不存在"))
		return
	}

	logger.InfoLogger.Printf("成功获取攻略详情，ID: %s", id)
	c.JSON(http.StatusOK, types.SuccessResponse(toGuideResponse(guide), "获取攻略成功"))
}

type SearchSuggestionResponse struct {
	Suggestions []string `json:"suggestions"`
}

// GetSearchSuggestions 获取搜索相关推荐词
func (gc *GuideController) GetSearchSuggestions(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusOK, types.SuccessResponse(
			SearchSuggestionResponse{Suggestions: []string{}},
			"获取推荐成功",
		))
		return
	}

	// 从标题和内容中搜索匹配的关键词
	var suggestions []string
	if err := gc.db.Model(&models.TravelGuide{}).
		Select("DISTINCT title").
		Where("title LIKE ?", "%"+keyword+"%").
		Limit(5).
		Pluck("title", &suggestions).Error; err != nil {
		c.JSON(http.StatusOK, types.ErrorResponse(1, "获取搜索推荐失败"))
		return
	}

	// 如果标题匹配不足5个，从内容中补充
	if len(suggestions) < 5 {
		var contentSuggestions []string
		if err := gc.db.Model(&models.TravelGuide{}).
			Select("DISTINCT title").
			Where("content LIKE ? AND title NOT IN ?", "%"+keyword+"%", suggestions).
			Limit(5-len(suggestions)).
			Pluck("title", &contentSuggestions).Error; err != nil {
			c.JSON(http.StatusOK, types.ErrorResponse(1, "获取搜索推荐失败"))
			return
		}
		suggestions = append(suggestions, contentSuggestions...)
	}

	c.JSON(http.StatusOK, types.SuccessResponse(
		SearchSuggestionResponse{Suggestions: suggestions},
		"获取推荐成功",
	))
}

// GetUserRecommendations 获取用户推荐攻略
func (gc *GuideController) GetUserRecommendations(c *gin.Context) {
	keyword := c.Query("keyword")
	offset, _ := c.GetQuery("offset")
	limit, _ := c.GetQuery("limit")

	// 设置默认值
	offsetInt := 0
	limitInt := 10
	if offset != "" {
		offsetInt, _ = strconv.Atoi(offset)
	}
	if limit != "" {
		limitInt, _ = strconv.Atoi(limit)
	}
	if limitInt <= 0 {
		limitInt = 10
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusOK, types.ErrorResponse(1, "未授权"))
		return
	}

	logger.InfoLogger.Printf("获取用户推荐 - 用户ID: %v, 关键词: %s, 偏移: %d, 限制: %d", userID, keyword, offsetInt, limitInt)

	// 构建基础查询
	query := gc.db.Model(&models.TravelGuide{}).
		Preload("User").
		Preload("Tags")

	// 检查用户是否有标签
	var userTagCount int64
	gc.db.Table("user_tags").Where("user_id = ?", userID).Count(&userTagCount)

	// 如果用户有标签，则按标签过滤
	if userTagCount > 0 {
		query = query.Joins("JOIN guide_tags ON guide_tags.guide_id = travel_guides.id").
			Joins("JOIN user_tags ON user_tags.tag_id = guide_tags.tag_id").
			Where("user_tags.user_id = ?", userID)
	}

	// 添加关键词搜索条件
	if keyword != "" {
		query = query.Where("(travel_guides.title LIKE ? OR travel_guides.content LIKE ?)", 
			"%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 查询数据
	var guides []models.TravelGuide
	query.Offset(offsetInt).Limit(limitInt + 1).Find(&guides) // 多查询一条用于判断是否还有更多

	// 处理分页结果
	hasMore := false
	if len(guides) > limitInt {
		hasMore = true
		guides = guides[:limitInt] // 去掉多查询的一条
	}

	// 转换响应格式
	guideResponses := make([]types.GuideResponse, 0, len(guides))
	for _, guide := range guides {
		guideResponses = append(guideResponses, toGuideResponse(guide))
	}

	c.JSON(http.StatusOK, types.SuccessResponse(
		gin.H{
			"list":    guideResponses,
			"hasMore": hasMore,
			"total":   total,
		},
		"获取推荐成功",
	))
}

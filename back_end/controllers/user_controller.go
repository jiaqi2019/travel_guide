package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"travel_guide/middleware"
	"travel_guide/models"
	"travel_guide/types"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

type CreateUserRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Nickname  string            `json:"nickname"`
	AvatarURL string            `json:"avatar_url"`
	Status    models.UserStatus `json:"status"`
}

type UpdateUserRoleRequest struct {
	Role models.UserRole `json:"role" binding:"required,oneof=admin user"`
}

type UpdateUserStatusRequest struct {
	Status models.UserStatus `json:"status" binding:"required,oneof=active banned"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type UserListResponse struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	Nickname   string    `json:"nickname"`
	AvatarURL  string    `json:"avatar_url"`
	Role       string    `json:"role"`
	Status     string    `json:"status"`
	GuideCount int64     `json:"guide_count"`
	CreatedAt  time.Time `json:"created_at"`
}

type GetUsersRequest struct {
	Limit  int `form:"limit,default=10"`
	Offset int `form:"offset,default=0"`
}

type PaginatedUserListResponse struct {
	List    []UserListResponse `json:"list"`
	Total   int64              `json:"total"`
	HasMore bool               `json:"has_more"`
}

func generateDefaultAvatar(nickname string) string {
	// 使用 DiceBear 的 avatars 风格生成头像
	// 使用用户名作为种子，确保每个用户有唯一的头像
	// 使用 SVG 格式，图片清晰且文件小
	// return fmt.Sprintf("https://api.dicebear.com/7.x/avatars/svg?seed=%s&backgroundType=gradient&backgroundColor=b6e3f4,c0aede,d1d4f9", nickname)
	// 使用昵称的第一个字符作为头像
	firstChar := string([]rune(nickname)[0])
	// 生成一个简单的默认头像URL
	// 这里使用一个示例URL，您可以根据需要替换为实际的默认头像服务
	return fmt.Sprintf("https://api.dicebear.com/7.x/initials/svg?seed=%s", firstChar)

}

func (uc *UserController) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, types.ErrorResponse(1, "请求参数错误"))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusOK, types.ErrorResponse(1, "密码加密失败"))
		return
	}

	// 如果用户没有提供头像URL，则生成默认头像
	avatarURL := req.AvatarURL
	if avatarURL == "" {
		avatarURL = generateDefaultAvatar(req.Nickname)
	}

	user := models.User{
		Username:  req.Username,
		Password:  string(hashedPassword),
		Nickname:  req.Nickname,
		AvatarURL: avatarURL,
		Role:      models.RoleUser,
		Status:    models.StatusActive,
	}

	if err := uc.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusOK, types.ErrorResponse(1, "创建用户失败"))
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse(
		gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"nickname":   user.Nickname,
			"avatar_url": user.AvatarURL,
			"role":       user.Role,
			"status":     user.Status,
			"created_at": user.CreatedAt,
		},
		"用户创建成功",
	))
}

func (uc *UserController) UpdateUserRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, types.ErrorResponse(1, "无效的用户ID"))
		return
	}

	var req UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, types.ErrorResponse(1, "请求参数错误"))
		return
	}

	var user models.User
	if err := uc.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusOK, types.ErrorResponse(1, "用户不存在"))
		return
	}

	if err := uc.DB.Model(&user).Update("role", req.Role).Error; err != nil {
		c.JSON(http.StatusOK, types.ErrorResponse(1, "更新用户角色失败"))
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse(
		gin.H{
			"id":   user.ID,
			"role": user.Role,
		},
		"更新用户角色成功",
	))
}

func (uc *UserController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, types.ErrorResponse(1, "请求参数错误"))
		return
	}

	var user models.User
	if err := uc.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, types.ErrorResponse(1, "用户不存在"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusOK, types.ErrorResponse(1, "密码错误"))
		return
	}

	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusOK, types.ErrorResponse(1, "生成Token失败"))
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse(
		gin.H{
			"token": token,
			"user": gin.H{
				"id":         user.ID,
				"username":   user.Username,
				"nickname":   user.Nickname,
				"avatar_url": user.AvatarURL,
				"role":       user.Role,
				"status":     user.Status,
			},
		},
		"登录成功",
	))
}

func (uc *UserController) GetUsers(c *gin.Context) {
	var req GetUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, types.ErrorResponse(1, "请求参数错误"))
		return
	}

	// 设置默认值
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	type UserWithCount struct {
		models.User
		GuideCount int64 `gorm:"column:guide_count"`
	}

	var users []UserWithCount
	var total int64

	// 获取总记录数
	if err := uc.DB.Model(&models.User{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusOK, types.ErrorResponse(1, "获取用户总数失败"))
		return
	}

	// 获取分页数据
	if err := uc.DB.Model(&models.User{}).
		Select("users.*, COUNT(guides.id) as guide_count").
		Joins("LEFT JOIN travel_guides as guides ON guides.user_id = users.id").
		Group("users.id").
		Offset(req.Offset).
		Limit(req.Limit + 1). // 多查询一条用于判断是否还有更多
		Find(&users).Error; err != nil {
		c.JSON(http.StatusOK, types.ErrorResponse(1, "获取用户列表失败"))
		return
	}

	// 判断是否还有更多数据
	hasMore := false
	if len(users) > req.Limit {
		hasMore = true
		users = users[:req.Limit] // 去掉多查询的一条
	}

	// 转换响应格式
	var userResponses []UserListResponse
	for _, user := range users {
		userResponses = append(userResponses, UserListResponse{
			ID:         user.ID,
			Username:   user.Username,
			Nickname:   user.Nickname,
			AvatarURL:  user.AvatarURL,
			Role:       string(user.Role),
			Status:     string(user.Status),
			GuideCount: user.GuideCount,
			CreatedAt:  user.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, types.SuccessResponse(
		PaginatedUserListResponse{
			List:    userResponses,
			Total:   total,
			HasMore: hasMore,
		},
		"获取用户列表成功",
	))
}

func (uc *UserController) UpdateUserStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, types.ErrorResponse(1, "无效的用户ID"))
		return
	}

	var req UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, types.ErrorResponse(1, "请求参数错误"))
		return
	}

	var user models.User
	if err := uc.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusOK, types.ErrorResponse(1, "用户不存在"))
		return
	}

	if err := uc.DB.Model(&user).Update("status", req.Status).Error; err != nil {
		c.JSON(http.StatusOK, types.ErrorResponse(1, "更新用户状态失败"))
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse(
		gin.H{
			"id":     user.ID,
			"status": user.Status,
		},
		"更新用户状态成功",
	))
}

package types

// 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 分页响应结构
type PaginationResponse struct {
	List    interface{} `json:"list"`
	Total   int64       `json:"total"`
	HasMore bool        `json:"has_more"`
}

// 成功响应
func SuccessResponse(data interface{}, message string) Response {
	if message == "" {
		message = "success"
	}
	return Response{
		Code:    0,
		Message: message,
		Data:    data,
	}
}

// 错误响应
func ErrorResponse(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}

// 各种响应数据结构
type GuideResponse struct {
	ID          uint          `json:"id"`
	Title       string        `json:"title"`
	Content     string        `json:"content"`
	Images      []string      `json:"images"`
	UserID      uint          `json:"user_id"`
	User        UserResponse  `json:"user"`
	PublishedAt int64         `json:"published_at"`
	Tags        []TagResponse `json:"tags"`
}

type CreateGuideResponse struct {
	ID          uint          `json:"id"`
	Title       string        `json:"title"`
	Content     string        `json:"content"`
	Images      []string      `json:"images"`
	PublishedAt int64         `json:"published_at"`
	Tags        []TagResponse `json:"tags"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
}

type TagResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

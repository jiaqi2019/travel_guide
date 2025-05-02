// 登录请求参数
export interface LoginRequest {
  username: string;
  password: string;
}

// 注册请求参数
export interface RegisterRequest {
  username: string;
  password: string;
  nickname: string;
}

// 登录响应
export interface LoginResponse {
  token: string;
  user: {
    id: number;
    username: string;
    nickname: string;
    role: 'admin' | 'user';
  };
}

// 注册响应
export interface RegisterResponse {
  message: string;
}

// API错误响应
export interface ApiError {
  message: string;
  code: number;
}

export interface Tag {
  id: number;
  name: string;
}

// 发布图文请求参数
export interface CreateGuideRequest {
  title: string;
  content: string;
  images: string[];
  tags: Tag[];
}

// 发布图文响应
export interface CreateGuideResponse {
  id: number;
  message: string;
}

// 图文列表项
export interface GuideItem {
  id: number;
  title: string;
  content: string;
  images: string[];
  tags: Tag[];
  user: {
    id: number;
    username: string;
    nickname: string;
    avatar: string;
  };
  published_at: number;
}

// 图文列表响应
export interface GuideListResponse {
  list: GuideItem[];
  has_more: boolean;
  total: number;
}

// 用户列表项
export interface UserItem {
  id: number;
  username: string;
  nickname: string;
  role: 'admin' | 'user';
  status: 'active' | 'banned';
  guide_count: number;
}

// 用户列表响应
export interface UserListResponse {
  list: UserItem[];
  total: number;
}

// 更新用户状态请求
export interface UpdateUserStatusRequest {
  user_id: number;
  status: 'active' | 'banned';
}

// 更新用户状态响应
export interface UpdateUserStatusResponse {
  message: string;
}

// 搜索建议响应
export interface SuggestionsResponse {
  suggestions: string[];
} 
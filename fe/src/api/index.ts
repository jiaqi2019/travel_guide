import axios, { AxiosInstance, InternalAxiosRequestConfig, AxiosResponse } from 'axios';
import { LoginRequest, RegisterRequest, LoginResponse, RegisterResponse, Tag, CreateGuideRequest, CreateGuideResponse, GuideListResponse, UserListResponse, UpdateUserStatusRequest, UpdateUserStatusResponse, SuggestionsResponse } from '../types/api';
import { ElMessage } from 'element-plus';

// 创建axios实例
const api: AxiosInstance = axios.create({
  baseURL: '/api', // 使用相对路径，由Vite代理处理
  headers: {
    'Content-Type': 'application/json'
  }
});

// 请求拦截器
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('token');
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
api.interceptors.response.use(
  (response: AxiosResponse) => {
    const { code, message, data } = response.data;
    if (code !== 0) {
      ElMessage.error(message || '请求失败');
      return Promise.reject(new Error(message));
    }
    return data;
  },
  (error) => {
    if (error.response?.status === 401) {
      // 未授权，清除token并跳转到登录页
      localStorage.removeItem('token');
      window.location.href = '/';
    }
    ElMessage.error(error.response?.data?.message || '请求失败');
    return Promise.reject(error);
  }
);

// 登录
export const login = (data: LoginRequest): Promise<LoginResponse> => {
  return api.post('/login', data);
};

// 注册
export const register = (data: RegisterRequest): Promise<RegisterResponse> => {
  return api.post('/register', data);
};

interface UploadResponse {
  url: string;
}

// 上传图片
export const uploadImage = (file: File): Promise<UploadResponse> => {
  const formData = new FormData();
  formData.append('image', file);
  return api.post('/upload/image', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  });
};

// 获取所有标签
export const getTags = (): Promise<Tag[]> => {
  return api.get('/tags');
};

// 发布图文
export const createGuide = (data: CreateGuideRequest): Promise<CreateGuideResponse> => {
  return api.post('/guides', data);
};

// 获取图文列表
export const getGuides = (offset: number = 0, tag?: string): Promise<GuideListResponse> => {
  return api.get('/guides', {
    params: {
      offset,
      tag
    }
  });
};

// 获取用户列表
export const getUserList = (): Promise<UserListResponse> => {
  return api.get('/users');
};

// 更新用户状态
export const updateUserStatus = (data: UpdateUserStatusRequest): Promise<UpdateUserStatusResponse> => {
  return api.put(`/users/${data.user_id}/status`, { status: data.status });
};

// 获取搜索建议
export const getSuggestions = (keyword: string): Promise<SuggestionsResponse> => {
  return api.get('/guides/suggestions', {
    params: {
      keyword
    }
  });
};

// 搜索图文
export const searchGuides = async (
  keyword: string,
  offset: number = 0,
  tag?: string
): Promise<GuideListResponse> => {
  return api.get('/guides/search', {
    params: { keyword, offset, tag }
  });
   
};

export const getRelatedTags = async (keyword: string): Promise<{ list: Tag[], has_more: boolean }> => {
  return await api.get('/tags/related', {
      params: { keyword }
  });
};

// 获取推荐攻略
export const getRecommendations = async (
  keyword?: string,
  offset: number = 0
): Promise<GuideListResponse> => {
  return api.get('/guides/recommendations', {
    params: { keyword, offset }
  });
};

export default api; 
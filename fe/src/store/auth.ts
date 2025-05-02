import { defineStore } from 'pinia';
import { ref } from 'vue';

interface UserInfo {
  id: number;
  username: string;
  nickname: string;
  avatar_url?: string;
  role: 'admin' | 'user';
  status: 'active' | 'banned';
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '');
  const isLoggedIn = ref(!!token.value);
  const userInfo = ref<UserInfo | null>(null);

  const setToken = (newToken: string) => {
    token.value = newToken;
    localStorage.setItem('token', newToken);
    isLoggedIn.value = true;
  };

  const setUserInfo = (info: UserInfo) => {
    userInfo.value = info;
    localStorage.setItem('userInfo', JSON.stringify(info));
  };

  const clearToken = () => {
    token.value = '';
    localStorage.removeItem('token');
    localStorage.removeItem('userInfo');
    isLoggedIn.value = false;
    userInfo.value = null;
  };

  // 初始化时从 localStorage 恢复用户信息
  const storedUserInfo = localStorage.getItem('userInfo');
  if (storedUserInfo) {
    userInfo.value = JSON.parse(storedUserInfo);
  }

  return {
    token,
    isLoggedIn,
    userInfo,
    setToken,
    setUserInfo,
    clearToken
  };
}); 
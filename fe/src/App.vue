<template>
  <div class="min-h-screen bg-white">
    <main class="container mx-auto flex  bg-white" style="height:calc(100vh-4rem);min-height:0;">
      <!-- 左侧边栏 -->
      <aside class="w-60 pr-8 flex flex-col gap-4">
        <div class="text-2xl font-bold text-red-500 w-48 h-[72px] items-center flex">旅游攻略</div>
        <el-button 
          :type="route.path === '/' ? 'primary' : 'default'" 
          size="large" 
          @click="goDiscover"
          class="!w-48 !ml-0"
        >
          发现
        </el-button>
        <el-button 
          :type="route.path === '/publish' ? 'primary' : 'default'" 
          size="large" 
          @click="goPublish"
          class="!w-48 !ml-0"
        >
          发布
        </el-button>
        <el-button 
          v-if="authStore.userInfo?.role === 'admin'"
          :type="route.path === '/users' ? 'primary' : 'default'" 
          size="large" 
          @click="goUsers"
          class="!w-48 !ml-0"
        >
          用户管理
        </el-button>
        <template v-if="authStore.isLoggedIn">
          <el-dropdown trigger="click" class="!w-48 !ml-0">
            <div class="flex items-center gap-2 px-4 py-2 cursor-pointer hover:bg-gray-100 rounded-lg w-full">
              <el-avatar :size="32" :src="authStore.userInfo?.avatar_url" />
              <span class="text-sm">{{ authStore.userInfo?.nickname }}</span>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <!-- <el-dropdown-item @click="router.push('/profile')">个人中心</el-dropdown-item> -->
                <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
        <el-button 
          v-else
          :type="route.path === '/login' ? 'primary' : 'default'" 
          size="large"
          class="!w-48 !ml-0"
          @click="showLoginDialog"
        >
          登录
        </el-button>
      </aside>
      <!-- 右侧攻略区域 -->
      <section class="flex-1 overflow-y-auto h-full min-h-0 hide-scrollbar">
        <div class="sticky top-0 h-[72px] flex justify-center items-center bg-[rgba(255,255,255,0.95)] backdrop-blur-sm z-20">
          <el-autocomplete
            v-model="searchKeyword"
            :fetch-suggestions="querySearch"
            placeholder="搜索旅游内容"
            class="w-96"
            @select="handleSuggestionSelect"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
            <template #suffix>
              <el-icon v-if="searchKeyword" class="cursor-pointer" @click="clearSearch"><Close /></el-icon>
            </template>
          </el-autocomplete>
        </div>
        <router-view v-slot="{ Component }">
          <component :is="Component" />
        </router-view>
      </section>
    </main>
    <login-dialog ref="loginDialogRef" />
  </div>
</template>

<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router';
import { ref } from 'vue';
import LoginDialog from './components/LoginDialog.vue';
import { useAuthStore } from './store/auth';
import { ElMessage } from 'element-plus';
import { Search, Close } from '@element-plus/icons-vue';
import { getSuggestions } from './api/index';

const router = useRouter();
const route = useRoute();
const loginDialogRef = ref<InstanceType<typeof LoginDialog>>();
const authStore = useAuthStore();
const searchKeyword = ref('');

const goDiscover = () => router.push('/');
const goPublish = () => {
  if (!authStore.isLoggedIn) {
    showLoginDialog();
    return;
  }
  router.push('/publish');
};
const goUsers = () => {
  if (!authStore.isLoggedIn) {
    showLoginDialog();
    return;
  }
  router.push('/users');
};
const showLoginDialog = () => {
  loginDialogRef.value?.show();
};

const handleLogout = () => {
  authStore.clearToken();
  ElMessage.success('已退出登录');
  router.push('/');
};

const querySearch = async (queryString: string, cb: (results: { value: string }[]) => void) => {
  if (!queryString) {
    cb([]);
    return;
  }
  
  try {
    const response = await getSuggestions(queryString);
    const results = response.suggestions.map(suggestion => ({
      value: suggestion
    }));
    cb(results);
  } catch (error) {
    console.error('Failed to fetch suggestions:', error);
    cb([]);
  }
};

const handleSuggestionSelect = (item) => {
  handleSearch();
};

const handleSearch = () => {
  if (searchKeyword.value.trim()) {
    router.push({
      path: '/search',
      query: { keyword: searchKeyword.value.trim() }
    });
  }
};

const clearSearch = () => {
  searchKeyword.value = '';
};
</script>

<style scoped>
main {
  height: calc(100vh);
  min-height: 0;
}
.hide-scrollbar {
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE 10+ */
}
.hide-scrollbar::-webkit-scrollbar {
  display: none; /* Chrome/Safari/Webkit */
}
.el-button {
  --el-button-hover-bg-color: transparent;
  --el-button-hover-border-color: var(--el-color-primary);
  --el-button-hover-text-color: var(--el-color-primary);
}
</style> 
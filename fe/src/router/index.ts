import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
import Feed from '../views/Feed.vue';
import PublishGuide from '../views/PublishGuide.vue';
import UserManagement from '../views/UserManagement.vue';
import SearchResults from '../views/SearchResults.vue';
import { useAuthStore } from '../store/auth';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: Feed,
  },
  {
    path: '/search',
    component: SearchResults,
  },
  {
    path: '/publish',
    component: PublishGuide,
    meta: { requiresAuth: true }
  },
  {
    path: '/users',
    component: UserManagement,
    meta: { requiresAuth: true, requiresAdmin: true }
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// 导航守卫
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();
  
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    // 如果需要登录但未登录，重定向到首页
    next('/');
  } else if (to.meta.requiresAdmin && authStore.userInfo?.role !== 'admin') {
    // 如果需要管理员权限但用户不是管理员，重定向到首页
    next('/');
  } else {
    next();
  }
});

export default router; 
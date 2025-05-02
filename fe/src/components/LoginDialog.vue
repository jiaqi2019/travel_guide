<template>
  <el-dialog
    v-model="dialogVisible"
    :title="isLogin ? '登录' : '注册'"
    width="400px"
    :close-on-click-modal="false"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="80px"
      class="login-form"
    >
      <el-form-item label="用户名" prop="username">
        <el-input v-model="form.username" placeholder="请输入用户名" />
      </el-form-item>
      <el-form-item label="密码" prop="password">
        <el-input
          v-model="form.password"
          type="password"
          placeholder="请输入密码"
          show-password
        />
      </el-form-item>
      <el-form-item v-if="!isLogin" label="昵称" prop="nickname">
        <el-input v-model="form.nickname" placeholder="请输入昵称" />
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="flex justify-between items-center">
        <el-button link @click="toggleMode">
          {{ isLogin ? '没有账号？去注册' : '已有账号？去登录' }}
        </el-button>
        <div class="flex gap-2">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="loading">确定</el-button>
        </div>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { ElMessage } from 'element-plus';
import { login, register } from '../api/index';
import type { FormInstance, FormRules } from 'element-plus';
import type { LoginRequest, RegisterRequest } from '../types/api';
import { useAuthStore } from '../store/auth';

const dialogVisible = ref(false);
const isLogin = ref(true);
const formRef = ref<FormInstance>();
const loading = ref(false);
const authStore = useAuthStore();

interface FormData {
  username: string;
  password: string;
  nickname: string;
}

const form = reactive<FormData>({
  username: '',
  password: '',
  nickname: '',
});

const rules = reactive<FormRules>({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ]
});

const toggleMode = () => {
  isLogin.value = !isLogin.value;
  formRef.value?.resetFields();
};

const handleSubmit = async () => {
  if (!formRef.value) return;
  
  try {
    await formRef.value.validate();
    loading.value = true;
    
    if (isLogin.value) {
      // 登录逻辑
      try {
        const data = await login({
          username: form.username,
          password: form.password
        } as LoginRequest);
        authStore.setToken(data.token);
        // 保存用户信息
        authStore.setUserInfo({
          id: data.user.id,
          username: data.user.username,
          nickname: data.user.nickname,
          avatar_url: data.user.avatar_url,
          role: data.user.role,
          status: data.user.status
        });
        ElMessage.success('登录成功');
        dialogVisible.value = false;
        // 刷新页面以更新状态
        // window.location.reload();
      } catch (error) {
        ElMessage.error((error as any).response?.data?.message || '登录失败');
      }
    } else {
      // 注册逻辑
      try {
        await register({
          username: form.username,
          password: form.password,
          nickname: form.nickname
        } as RegisterRequest);
        ElMessage.success('注册成功');
        isLogin.value = true;
        formRef.value?.resetFields();
      } catch (error) {
        ElMessage.error((error as any).response?.data?.message || '注册失败');
      }
    }
  } catch (error) {
    console.error('表单验证失败:', error);
  } finally {
    loading.value = false;
  }
};

// 暴露方法给父组件
defineExpose({
  show: () => {
    dialogVisible.value = true;
  }
});
</script>

<style scoped>
.login-form {
  padding: 20px 0;
}
</style> 
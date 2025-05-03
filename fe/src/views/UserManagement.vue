<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold mb-6">用户管理</h1>
    <el-table :data="users" v-loading="loading" style="width: 100%">
      <el-table-column prop="id" label="ID"  />
      <el-table-column prop="username" label="用户名"  />
      <el-table-column prop="nickname" label="昵称"  />
      <el-table-column prop="guide_count" label="发文数" />
      <el-table-column label="兴趣标签">
        <template #default="{ row }">
          <el-tag
            v-for="tag in row.tags"
            :key="tag"
            class="mr-1"
            size="small"
            type="warning"
          >
            {{ tag }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="role" label="角色" >
        <template #default="{ row }">
          <el-tag :type="row.role === 'admin' ? 'danger' : 'info'">
            {{ row.role === 'admin' ? '管理员' : '普通用户' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" >
        <template #default="{ row }">
          <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
            {{ row.status === 'active' ? '正常' : '已禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" >
        <template #default="{ row }">
          <el-button
            v-if="row.role !== 'admin'"
            :type="row.status === 'active' ? 'danger' : 'success'"
            size="small"
            @click="toggleUserStatus(row)"
          >
            {{ row.status === 'active' ? '禁用' : '启用' }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { ElMessage } from 'element-plus';
import { getUserList, updateUserStatus } from '../api/index';
import type { UserItem } from '../types/api';

const users = ref<UserItem[]>([]);
const loading = ref(false);

// 获取用户列表
const fetchUsers = async () => {
  try {
    loading.value = true;
    const response = await getUserList();
    users.value = response.list;
  } catch (error) {
    ElMessage.error('获取用户列表失败');
  } finally {
    loading.value = false;
  }
};

// 切换用户状态
const toggleUserStatus = async (user: UserItem) => {
  try {
    await updateUserStatus({
      user_id: user.id,
      status: user.status === 'active' ? 'banned' : 'active'
    });
    // 直接修改本地数据
    const index = users.value.findIndex(u => u.id === user.id);
    if (index !== -1) {
      users.value[index] = {
        ...users.value[index],
        status: user.status === 'active' ? 'banned' : 'active'
      };
    }
    ElMessage.success(`已${user.status === 'active' ? '禁用' : '启用'}用户`);
  } catch (error) {
    ElMessage.error('操作失败');
  }
};

onMounted(() => {
  fetchUsers();
});
</script>

<style scoped>
.el-table {
  border-radius: 4px;
  overflow: hidden;
}

/* 确保表格容器不会超出预期宽度 */
.p-6 {
  max-width: 100%;
  overflow-x: auto;
}
</style> 
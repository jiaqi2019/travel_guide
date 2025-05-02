<template>
  <div class="w-full bg-white">
    <div class="h-[72px] flex items-center px-4 border-b">
      <el-button link @click="goBack" class="mr-2 text-lg">&lt;</el-button>
      <span class="text-xl font-bold">发布图文</span>
      <el-button link class="ml-auto text-red-400" @click="clearImages">清空并重新上传</el-button>
    </div>
    <div class="p-4">
      <div class="mb-6">
        <div class="font-semibold mb-2">图片编辑 <span class="text-gray-400 text-sm">({{ form.images.length }}/9 )</span></div>
        <el-upload
          class="upload-demo"
          action="#"
          list-type="picture-card"
          :auto-upload="false"
          :on-preview="handlePreview"
          :on-remove="handleRemove"
          :file-list="form.images"
          :limit="9"
          :on-exceed="handleExceed"
          :on-change="handleChange"
          multiple
        >
          <el-icon><Plus /></el-icon>
        </el-upload>
        <el-dialog v-model="dialogVisible" title="图片预览">
          <img w-full :src="dialogImageUrl" alt="Preview Image" />
        </el-dialog>
      </div>
      <div class="mb-6">
        <div class="font-semibold mb-2">正文内容</div>
        <el-input
          v-model="form.title"
          maxlength="20"
          show-word-limit
          placeholder="填写标题会有更多赞哦 ~"
          class="mb-2"
        />
        <el-input
          v-model="form.content"
          type="textarea"
          :rows="6"
          maxlength="1000"
          show-word-limit
          placeholder="输入正文描述，真诚有价值的分享令人温暖"
        />
      </div>
      <div class="flex items-center gap-2 pt-4">
        <el-select v-model="form.topic" placeholder="# 话题" class="w-32">
          <el-option
            v-for="tag in tags"
            :key="tag.id"
            :label="tag.name"
            :value="tag.name"
          />
        </el-select>
        <el-button type="primary" class="ml-auto px-8" @click="onSubmit">发布</el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { Plus } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';
import { uploadImage, getTags, createGuide } from '../api/index';

const router = useRouter();
const dialogVisible = ref(false);
const dialogImageUrl = ref('');
const tags = ref([]);

const form = ref({
  title: '',
  content: '',
  images: [],
  topic: '',
  user: '',
});

// 获取所有标签
const fetchTags = async () => {
  try {
    const response = await getTags();
    tags.value = response;
  } catch (error) {
    ElMessage.error('获取标签失败');
  }
};

onMounted(() => {
  fetchTags();
});

const handlePreview = (file) => {
  dialogImageUrl.value = file.url;
  dialogVisible.value = true;
};

const handleRemove = (file, fileList) => {
  form.value.images = fileList;
};

// 在 handleChange 函数中：
const handleChange = async (file, fileList) => {
  try {
    const response = await uploadImage(file.raw);
    const uploadedFile = {
      ...file,
      url: response.url
    };
    console.log('uploadedFile', uploadedFile);
    form.value.images = fileList.map(f => 
      f.uid === file.uid ? uploadedFile : f
    );
  } catch (error) {
    console.error('Upload failed:', error);
    ElMessage.error('图片上传失败');
    form.value.images = fileList.filter(f => f.uid !== file.uid);
  }
  // form.value.images = fileList;
};

const handleExceed = () => {
  ElMessage.warning('最多只能上传9张图片');
};

const clearImages = () => {
  form.value.images = [];
};

const onSubmit = async () => {
  try {
    // 表单验证
    if (!form.value.title) {
      ElMessage.warning('请输入标题');
      return;
    }
    if (!form.value.content) {
      ElMessage.warning('请输入正文内容');
      return;
    }
    if (form.value.images.length === 0) {
      ElMessage.warning('请至少上传一张图片');
      return;
    }
    if (!form.value.topic) {
      ElMessage.warning('请选择话题');
      return;
    }

    // 准备发布数据
    const guideData = {
      title: form.value.title,
      content: form.value.content,
      images: form.value.images.map(img => img.url),
      tags: [form.value.topic]
    };

    // 调用发布接口
    await createGuide(guideData);
    ElMessage.success('发布成功！');
    router.push('/');
  } catch (error) {
    console.error('发布失败:', error);
  }
};

const goBack = () => {
  router.back();
};
</script>

<style scoped>
.upload-demo .el-upload {
  display: flex;
  flex-wrap: wrap;
}
</style> 
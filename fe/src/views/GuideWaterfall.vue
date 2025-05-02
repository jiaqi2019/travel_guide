<template>
  <div class="w-full">
    <GuideWaterfallList
      :tag-data="tagData"
      :active-tag="activeTag"
      :loading="loading"
      :loading-more="loadingMore"
      :show-tags="true"
      :empty-title="'暂无内容'"
      :empty-subtitle="'尝试其他操作'"
      :empty-description="'暂无内容'"
      @load-more="fetchGuides"
      @tag-click="handleTagClick"
    />

    <!-- 自定义弹窗 -->
    <div v-if="dialogVisible" class="custom-modal">
      <div class="modal-overlay" @click="dialogVisible = false"></div>
      <div class="modal-content">
        <div class="modal-header">
          <div class="modal-title">{{ currentGuide?.title }}</div>
          <div class="modal-close" @click="dialogVisible = false">
            <el-icon><Close /></el-icon>
          </div>
        </div>
        <div class="modal-body">
          <div class="guide-detail-content">
            <!-- 左侧图片轮播 -->
            <div class="guide-images">
              <el-carousel 
                height="100%" 
                :autoplay="false"
                trigger="click"
                indicator-position="outside"
                class="rounded-lg overflow-hidden h-full"
              >
                <el-carousel-item v-for="(image, index) in currentGuide?.images" :key="index" class="h-full">
                  <img 
                    :src="formatImageUrl(image)" 
                    class="max-h-full max-w-full object-contain"
                    alt=""
                  />
                </el-carousel-item>
              </el-carousel>
            </div>
            <!-- 右侧内容区 -->
            <div class="guide-info">
              <div class="user-info">
                <div class="flex items-center gap-4">
                  <img 
                    :src="formatImageUrl(currentGuide?.user?.avatar_url)" 
                    class="w-12 h-12 rounded-full object-cover"
                    alt=""
                  />
                  <div>
                    <div class="font-semibold text-lg">{{ currentGuide?.user?.nickname }}</div>
                    <div class="text-gray-500 text-sm">{{ formatDate(currentGuide?.published_at) }}</div>
                  </div>
                </div>
                <div class="tags mt-4">
                  <el-tag
                    v-for="tag in currentGuide?.tags"
                    :key="tag.id"
                    size="small"
                    effect="plain"
                    class="mr-2 mb-2"
                  >
                    {{ tag.name }}
                  </el-tag>
                </div>
              </div>
              <div class="guide-content">
                <div class="text-gray-700 whitespace-pre-wrap">{{ currentGuide?.content }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import GuideWaterfallList from '../components/GuideWaterfallList.vue';
import { getGuides, getTags } from '../api/index';
import { Loading, Close } from '@element-plus/icons-vue'

const activeTag = ref('全部');
const loading = ref(true);
const loadingMore = ref(false);

// 为每个tag维护独立的数据
const tagData = ref({
  '全部': {
    guides: [],
    offset: 0,
    hasMore: true
  }
});

const fetchTags = async () => {
  try {
    const response = await getTags();
    // 初始化每个tag的数据
    response.forEach(tag => {
      if (!tagData.value[tag.name]) {
        tagData.value[tag.name] = {
          guides: [],
          offset: 0,
          hasMore: true
        };
      }
    });
  } catch (error) {
    console.error('Failed to fetch tags:', error);
  }
};

const fetchGuides = async (tag: string, reset = false) => {
  if (reset) {
    loading.value = true;
    tagData.value[tag].offset = 0;
  } else {
    loadingMore.value = true;
  }
  
  try {
    const {list, has_more} = await getGuides(
      tagData.value[tag].offset,
      tag === '全部' ? undefined : tag
    );
    
    if (reset) {
      tagData.value[tag].guides = list;
    } else {
      tagData.value[tag].guides = [...tagData.value[tag].guides, ...list];
    }
    tagData.value[tag].hasMore = has_more;
    tagData.value[tag].offset = tagData.value[tag].guides.length;
  } catch (error) {
    console.error('Failed to fetch guides:', error);
  } finally {
    loading.value = false;
    loadingMore.value = false;
  }
};

const handleTagClick = (tagName: string) => {
  activeTag.value = tagName;
  // 如果该tag的数据为空，则加载数据
  if (tagData.value[tagName].guides.length === 0) {
    fetchGuides(tagName, true);
  }
};

const dialogVisible = ref(false);
const currentGuide = ref(null);

const showDetail = (guide) => {
  currentGuide.value = guide;
  dialogVisible.value = true;
};

const formatDate = (timestamp) => {
  if (!timestamp) return '';
  const date = new Date(timestamp * 1000);
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

onMounted(() => {
  fetchTags();
  fetchGuides('全部', true);
});
</script>

<style scoped>
.waterfall-container :deep(.waterfall-item) {
  width: 240px !important;
}

/* 添加过渡动画 */
.group {
  transition: transform 0.3s ease;
}

.group:hover {
  transform: translateY(-4px);
}

.custom-modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1000;
}

.modal-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
}

.modal-content {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 60vw;
  height: 90vh;
  background-color: white;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e4e7ed;
}

.modal-title {
  font-size: 18px;
  font-weight: 600;
}

.modal-close {
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.modal-close:hover {
  background-color: #f5f7fa;
}

.modal-body {
  flex: 1;
  overflow: hidden;
  padding: 20px;
}

.guide-detail-content {
  display: flex;
  gap: 2rem;
  height: 100%;
}

.guide-images {
  flex: 1;
  border-right: 1px solid #e4e7ed;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.guide-info {
  max-width: 440px;
  min-width: 300px;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.user-info {
  padding-bottom: 20px;
  border-bottom: 1px solid #e4e7ed;
  flex-shrink: 0;
}

.guide-content {
  flex: 1;
  overflow-y: auto;
  padding-top: 20px;
  padding-right: 8px;
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE and Edge */
}

.guide-content::-webkit-scrollbar {
  display: none; /* Chrome, Safari, Opera */
}

.guide-content::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.guide-content::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 3px;
}

.guide-content::-webkit-scrollbar-thumb:hover {
  background: #555;
}

:deep(.el-carousel) {
  width: 100%;
  height: 100%;
}

:deep(.el-carousel__container) {
  height: 100%;
}

:deep(.el-carousel__item) {
  display: flex;
  align-items: center;
  justify-content: center;
}

:deep(.el-carousel__indicators) {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 10;
}

:deep(.el-carousel__button) {
  width: 30px;
  height: 4px;
  border-radius: 2px;
  background-color: rgba(0, 0, 0, 0.3);
  transition: all 0.3s;
}

:deep(.el-carousel__indicator.is-active .el-carousel__button) {
  background-color: #409EFF;
  width: 40px;
}

/* 响应式布局 */
@media screen and (max-width: 1446px) {
  .modal-content {
    width: 90vw;
    max-width: none;
  }

  .guide-detail-content {
    flex-direction: column;
    gap: 1rem;
  }

  .guide-images {
    width: 100%;
    min-width: auto;
    border-right: none;
    border-bottom: 1px solid #e4e7ed;
    padding-bottom: 1rem;
    height: 400px;
  }

  .guide-info {
    width: 100%;
    max-width: none;
    min-width: auto;
    height: auto;
  }
}
</style> 
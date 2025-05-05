<template>
  <div class="w-full">
    <!-- 地点tag tab -->
    <div v-if="showTags" class="flex items-center h-[72px] gap-4 mb-6 sticky top-[72px] z-10 bg-[rgba(255,255,255,0.95)] backdrop-blur-sm">
      <el-tag
        v-for="tagName in Object.keys(tagData)"
        :key="tagName"
        :type="tagName === activeTag ? 'success' : 'info'"
        class="cursor-pointer text-base px-4 py-2"
        @click="handleTagClick(tagName)"
        effect="plain"
      >
        {{ tagName }}
      </el-tag>
    </div>
    <div v-loading="loading" v-infinite-scroll="loadMore" :infinite-scroll-disabled="loadingMore || !currentHasMore">
      <template v-if="!loading && currentGuides.length === 0">
        <div class="flex flex-col items-center justify-center py-12">
          <el-empty
            :image-size="200"
            :description="emptyDescription"
          >
            <template #description>
              <div class="text-gray-500 text-lg">{{ emptyTitle }}</div>
              <div class="text-gray-400 text-sm mt-2">{{ emptySubtitle }}</div>
            </template>
          </el-empty>
        </div>
      </template>
      <template v-else>
        <Waterfall
          :list="currentGuides"
          :gutter="16"
          :colNum="5"
          :width="265"
          :lazyload="true"
          :crossOrigin="true"
        >
          <template #item="{ item }">
            <div class="relative group" @click="showDetail(item)">
              <LazyImg
                :url="item.images[0]" 
                class="w-full h-48 object-cover cursor-pointer rounded-[16px] transition-all duration-300 ease-linear"
                @load="imageLoad"
                @error="imageError"
                @success="imageSuccess"
              />
              <div class=" cursor-pointer absolute inset-0 bg-black bg-opacity-40 opacity-0  group-hover:opacity-100 transition-all duration-300 rounded-[16px]"></div>
            </div>
            <div class="p-3">
              <div class="text-gray-400 text-[14px] mb-2">{{ item.title }}</div>
              <div class="flex items-center gap-2 mt-2">
                <img :src="(item.user.avatar_url)" class="w-5 h-5 rounded-full object-cover" />
                <span class="text-gray-400 text-xs">{{ item.user.nickname }}</span>
              </div>
            </div>
          </template>
        </Waterfall>
        <div v-if="loadingMore" class="text-center py-4">
          <el-icon class="is-loading"><Loading /></el-icon>
          <span class="ml-2">加载中...</span>
        </div>
      </template>
    </div>

    <!-- 自定义弹窗 -->
    <div v-if="dialogVisible" class="custom-modal">
      <div class="modal-overlay" @click="dialogVisible = false">
        <div class="modal-close" @click.stop="dialogVisible = false">
          <el-icon><Close /></el-icon>
        </div>
      </div>
      <div class="modal-content">
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
                    :src="image" 
                    class="h-full w-full object-cover"
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
                    :src="currentGuide?.user?.avatar_url" 
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
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
import { LazyImg, Waterfall } from 'vue-waterfall-plugin-next'
import 'vue-waterfall-plugin-next/dist/style.css'
import { Loading, Close,  } from '@element-plus/icons-vue'
import type { GuideItem,  } from '../types/api';

interface Props {
  tagData: Record<string, {
    guides: GuideItem[];
    offset: number;
    hasMore: boolean;
  }>;
  activeTag: string;
  loading: boolean;
  loadingMore: boolean;
  showTags?: boolean;
  emptyTitle?: string;
  emptySubtitle?: string;
  emptyDescription?: string;
}

const props = withDefaults(defineProps<Props>(), {
  tagData: () => ({}),
  activeTag: '',
  loading: false,
  loadingMore: false,
  showTags: false,
  emptyTitle: '暂无内容',
  emptySubtitle: '尝试其他操作',
  emptyDescription: '暂无内容'
});

const emit = defineEmits<{
  (e: 'loadMore', tag: string): void;
  (e: 'tagClick', tag: string): void;
}>();

const dialogVisible = ref(false);
const currentGuide = ref<GuideItem | null>(null);

const currentGuides = computed(() => {
  return props.tagData[props.activeTag]?.guides || [];
});

const currentHasMore = computed(() => {
  const hasMore = props.tagData[props.activeTag]?.hasMore || false;
  console.log('currentHasMore', {
    activeTag: props.activeTag,
    hasMore,
    tagData: props.tagData[props.activeTag]
  });
  return hasMore;
});

const imageLoad = () => {
  console.log('image loaded');
};

const imageError = () => {
  console.error('image load failed');
};

const imageSuccess = () => {
  console.log('image load success');
};

const handleTagClick = (tagName: string) => {
  emit('tagClick', tagName);
};

const showDetail = (guide: GuideItem) => {
  currentGuide.value = guide;
  dialogVisible.value = true;
};

const formatDate = (timestamp: number) => {
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

const loadMore = () => {
  if (!props.loadingMore && currentHasMore.value) {
    emit('loadMore', props.activeTag);
  }
};

// 添加防抖函数
const debounce = (fn: Function, delay: number) => {
  let timer: number | null = null;
  return (...args: any[]) => {
    if (timer) {
      clearTimeout(timer);
    }
    timer = setTimeout(() => {
      fn(...args);
      timer = null;
    }, delay);
  };
};

// 使用防抖的加载更多函数
const debouncedLoadMore = debounce(loadMore, 300);

// 监听activeTag变化，重置加载状态
watch(() => props.activeTag, (newTag, oldTag) => {
  if (newTag !== oldTag) {
    // 只在标签真正变化时触发加载
    debouncedLoadMore();
  }
});

</script>

<style scoped>
.waterfall-container :deep(.waterfall-item) {
  width: 240px !important;
}

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

.modal-close {
  position: absolute;
  top: 20px;
  right: 20px;
  cursor: pointer;
  padding: 8px;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.1);
  transition: all 0.3s;
  z-index: 1001;
  color: white;
}

.modal-close:hover {
  background-color: rgba(255, 255, 255, 0.2);
}

.modal-close :deep(.el-icon) {
  color: white;
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
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.guide-content::-webkit-scrollbar {
  display: none;
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
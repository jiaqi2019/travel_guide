<template>
  <div class="w-full">
    <GuideWaterfallList
      :tag-data="tagData"
      :active-tag="activeTag"
      :loading="loading"
      :loading-more="loadingMore"
      :show-tags="true"
      :empty-title="'没有找到相关的内容'"
      :empty-subtitle="'尝试使用其他关键词搜索'"
      :empty-description="'暂无相关搜索结果'"
      @load-more="fetchGuides"
      @tag-click="handleTagClick"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import GuideWaterfallList from '../components/GuideWaterfallList.vue';
import { searchGuides, getRelatedTags } from '../api/index';

const route = useRoute();
const keyword = computed(() => route.query.keyword);
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
    const response = await getRelatedTags(keyword.value);
    // 初始化每个tag的数据
    response.list.forEach(tag => {
      if (!tagData.value[tag.name]) {
        tagData.value[tag.name] = {
          guides: [],
          offset: 0,
          hasMore: true
        };
      }
    });
  } catch (error) {
    console.error('Failed to fetch related tags:', error);
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
    const {list, has_more} = await searchGuides(
      keyword.value,
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
    console.error('Failed to fetch search results:', error);
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

// 监听 keyword 变化
watch(keyword, (newKeyword) => {
  if (newKeyword) {
    // 重置状态
    activeTag.value = '全部';
    tagData.value = {
      '全部': {
        guides: [],
        offset: 0,
        hasMore: true
      }
    };
    // 重新获取数据
    fetchTags();
    fetchGuides('全部', true);
  }
}, { immediate: true });

onMounted(() => {
  fetchTags();
  fetchGuides('全部', true);
});

</script>

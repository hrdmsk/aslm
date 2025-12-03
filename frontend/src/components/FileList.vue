<template>
  <div class="file-list-container">
    <!-- コンテンツエリア：リストとプレビューを横並びにする -->
    <div class="content-area">
      
      <!-- 左側：ファイル一覧 (グリッド表示) -->
      <div class="list-body grid-view">
        <div 
          v-for="item in store.fileList" 
          :key="item.name"
          class="grid-item"
          @dblclick="handleDoubleClick(item)"
          :class="{ 'is-selected': selectedItemName === item.name }"
          @click="selectItem(item)"
        >
          <!-- サムネイル画像エリア -->
          <div class="thumbnail-wrapper">
            <!-- フォルダの場合: イメージ画像を表示 -->
            <img 
              v-if="item.type === 'folder'"
              src="https://go.dev/images/gophers/motorcycle.svg" 
              alt="Folder Thumbnail" 
              class="thumbnail-image"
            />
            
            <!-- ファイルの場合: 大きなアイコンを表示 -->
            <span v-else class="file-icon-large">📄</span>
          </div>

          <!-- アイテム詳細 -->
          <div class="item-details">
            <div class="item-name" :title="item.name">
              {{ item.name }}
            </div>
            <div class="item-meta">
              <span v-if="item.source === 'booth'" class="tag booth" title="Booth Linked">B</span>
              <span v-else-if="item.source === 'gumroad'" class="tag gumroad" title="Gumroad Linked">G</span>
              <span class="item-date">2025/11/27</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 右側：プレビューペイン (別コンポーネント化) -->
      <FilePreview 
        :item="selectedItemData" 
      />

    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useFileSystemStore } from '../stores/fileSystem';
import FilePreview from './FilePreview.vue'; // インポート

const store = useFileSystemStore();
const selectedItemName = ref(null);

// 選択されたアイテムのオブジェクトデータを取得
const selectedItemData = computed(() => {
  return store.fileList.find(item => item.name === selectedItemName.value);
});

// アイテム選択処理
const selectItem = (item) => {
  selectedItemName.value = item.name;
};

// ダブルクリック時の挙動
const handleDoubleClick = (item) => {
  if (item.type === 'folder') {
    store.changeDirectory(`${store.currentPath}/${item.name}`);
    selectedItemName.value = null; // ディレクトリ移動したら選択解除
  }
};
</script>

<style scoped>
.file-list-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #ffffff;
  font-family: 'Segoe UI', sans-serif;
  font-size: 13px;
  overflow: hidden; /* コンテナ自体はスクロールさせない */
}

/* 左右分割用のラッパー */
.content-area {
  display: flex;
  flex: 1;
  overflow: hidden;
}

/* --- 左側：リストビュー --- */
.list-body.grid-view {
  flex: 1;
  overflow-y: auto;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 12px;
  padding: 12px;
  align-content: start;
}

.grid-item {
  display: flex;
  flex-direction: column;
  cursor: default;
  border: 1px solid transparent;
  border-radius: 4px;
  padding: 6px;
  text-align: center;
  transition: all 0.1s ease;
}

.grid-item:hover {
  background-color: #f0f9ff;
  border-color: #e5f3ff;
}

.grid-item.is-selected {
  background-color: #cce8ff;
  border-color: #99d1ff;
}

/* サムネイル (リスト内) */
.thumbnail-wrapper {
  width: 100%;
  aspect-ratio: 1 / 1;
  margin-bottom: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  font-size: 48px;
}

.thumbnail-image {
  width: auto;
  height: auto;
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.file-icon-large {
  opacity: 0.6;
  color: #888;
}

.item-details {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.item-name {
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.item-meta {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-size: 11px;
  color: #666;
}

/* ソース連携タグ */
.tag {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 10px;
  font-weight: bold;
}
.tag.booth { background-color: #fff2f2; color: #fc4d50; border: 1px solid #ffd1d1; }
.tag.gumroad { background-color: #effcf6; color: #26a17b; border: 1px solid #bcebdc; }
</style>
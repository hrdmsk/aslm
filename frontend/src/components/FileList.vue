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
            <!-- フォルダの場合: イメージ画像があれば表示、なければデフォルトアイコン -->
            <img 
              v-if="item.type === 'folder' && item.imageUrl"
              :src="item.imageUrl" 
              alt="Folder Thumbnail" 
              class="thumbnail-image"
            />
            <div v-else-if="item.type === 'folder'" class="folder-icon-large">📁</div>
            
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
  background-color: #f8fafc;
  font-family: 'Inter', 'Segoe UI', sans-serif;
  font-size: 14px;
  overflow: hidden;
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
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 16px;
  padding: 24px;
  align-content: start;
}

.grid-item {
  display: flex;
  flex-direction: column;
  cursor: pointer;
  background-color: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 12px;
  text-align: center;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}

.grid-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  border-color: #cbd5e1;
}

.grid-item.is-selected {
  border-color: #6366f1;
  background-color: #eef2ff;
  box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.2);
}

/* サムネイル (リスト内) */
.thumbnail-wrapper {
  width: 100%;
  aspect-ratio: 1 / 1;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  font-size: 48px;
  border-radius: 8px;
  background-color: #f8fafc;
}

.thumbnail-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.file-icon-large,
.folder-icon-large {
  opacity: 0.8;
  color: #94a3b8;
}

.folder-icon-large {
  color: #fbbf24; /* Amber 400 */
}

.item-details {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.item-name {
  font-weight: 600;
  color: #1e293b;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-size: 13px;
}

.item-meta {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-size: 11px;
  color: #64748b;
}

/* ソース連携タグ */
.tag {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 2px 8px;
  border-radius: 9999px;
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 0.025em;
}
.tag.booth { background-color: #fef2f2; color: #ef4444; border: 1px solid #fee2e2; }
.tag.gumroad { background-color: #ecfdf5; color: #10b981; border: 1px solid #d1fae5; }
</style>
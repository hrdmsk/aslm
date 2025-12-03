<template>
  <aside class="preview-pane">
    <template v-if="item">
      <div class="preview-header">
        <span class="preview-icon-small">{{ item.type === 'folder' ? '📁' : '📄' }}</span>
        <span class="preview-title">{{ item.name }}</span>
      </div>

      <div class="preview-image-container">
        <!-- フォルダの場合: 大きなプレビュー画像 -->
        <img 
          v-if="item.type === 'folder'"
          src="https://go.dev/images/gophers/motorcycle.svg" 
          alt="Preview" 
          class="preview-image-large"
        />
        <!-- ファイルの場合: 大きなアイコン -->
        <span v-else class="preview-file-icon">📄</span>
      </div>

      <div class="preview-info">
        <div class="info-row">
          <label>種類:</label>
          <span>{{ item.type === 'folder' ? 'ファイル フォルダー' : 'Unity Package' }}</span>
        </div>
        <div class="info-row">
          <label>更新日時:</label>
          <span>2025/11/27 18:00</span>
        </div>
        <div class="info-row" v-if="item.source">
          <label>Source:</label>
          <span :class="['tag', item.source]">
            {{ item.source === 'booth' ? 'Booth' : 'Gumroad' }}
          </span>
        </div>
        <div class="info-row" v-if="item.url">
          <label>Link:</label>
          <a :href="item.url" target="_blank" class="link-text">商品ページを開く</a>
        </div>
      </div>
    </template>
    <div v-else class="empty-state">
      <p>アイテムを選択してください</p>
    </div>
  </aside>
</template>

<script setup>
defineProps({
  item: {
    type: Object,
    required: false,
    default: null
  }
});
</script>

<style scoped>
.preview-pane {
  width: 280px; /* 固定幅 */
  border-left: 1px solid #e5e5e5;
  background-color: #fcfcfc;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  overflow-y: auto;
  flex-shrink: 0;
}

.preview-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #333;
  word-break: break-all;
}

.preview-image-container {
  width: 100%;
  aspect-ratio: 16 / 9;
  background-color: #fff;
  border: 1px solid #e5e5e5;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border-radius: 4px;
}

.preview-image-large {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.preview-file-icon {
  font-size: 64px;
  opacity: 0.5;
}

.preview-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
  font-size: 12px;
}

.info-row {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.info-row label {
  color: #666;
  font-size: 11px;
}

.info-row span {
  color: #333;
}

.link-text {
  color: #0066cc;
  text-decoration: none;
}
.link-text:hover {
  text-decoration: underline;
}

/* Tag Styles */
.tag {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 10px;
  font-weight: bold;
  width: fit-content;
}
.tag.booth { background-color: #fff2f2; color: #fc4d50; border: 1px solid #ffd1d1; }
.tag.gumroad { background-color: #effcf6; color: #26a17b; border: 1px solid #bcebdc; }

.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #999;
  font-size: 14px;
}
</style>
<template>
  <div class="file-preview" v-if="item">
    <!-- ヘッダー：アイコンと名前 -->
    <div class="preview-header">
      <div class="preview-icon">
        <span v-if="item.type === 'folder'">📁</span>
        <span v-else>📄</span>
      </div>
      <div class="preview-title">
        <h2>{{ item.name }}</h2>
        <p class="file-type">{{ item.type === 'folder' ? 'フォルダ' : 'ファイル' }}</p>
      </div>
    </div>

    <!-- プレビュー画像 -->
    <div class="preview-image-container">
      <img v-if="item.imageUrl" :src="item.imageUrl" alt="Preview" class="preview-image" />
      <div v-else class="no-image">No Image</div>
    </div>

    <!-- 閲覧モード -->
    <div class="preview-details" v-if="!isEditing">
      <div class="detail-item" v-if="item.url">
        <label>Link</label>
        <a :href="item.url" target="_blank" class="link-text">商品ページを開く</a>
      </div>

      <div class="detail-item" v-if="item.tags && item.tags.length > 0">
        <label>Tags</label>
        <div class="tags-list">
          <span v-for="tag in item.tags" :key="tag" class="tag">{{ tag }}</span>
        </div>
      </div>

      <div class="detail-item" v-if="item.path">
        <label>Path</label>
        <div class="path-text">{{ item.path }}</div>
      </div>

      <div class="actions">
        <button class="edit-btn" @click="startEditing">編集 ✎</button>
      </div>
    </div>

    <!-- 編集モード -->
    <div class="preview-details edit-mode" v-else>
      <div class="detail-item">
        <label>Booth/Product URL</label>
        <div class="input-with-button">
          <input type="text" v-model="editUrl" placeholder="https://booth.pm/..." />
          <button v-if="editUrl" @click="openUrl(editUrl)" title="開く">🔗</button>
        </div>
      </div>

      <div class="detail-item">
        <label>Thumbnail URL</label>
        <input type="text" v-model="editImageUrl" placeholder="https://..." />
      </div>

      <div class="detail-item">
        <label>Tags (comma separated)</label>
        <input type="text" v-model="editTags" placeholder="avatar, cute, 3d" />
      </div>

      <div class="actions">
        <button class="cancel-btn" @click="cancelEditing">キャンセル</button>
        <button class="save-btn" @click="saveChanges">保存</button>
      </div>
    </div>
  </div>
  <div v-else class="empty-state">
    アイテムを選択してください
  </div>
</template>

<script setup>
import { ref, watch } from 'vue';
import { UpdateProduct } from '../../wailsjs/go/main/App';
import { BrowserOpenURL } from '../../wailsjs/runtime';

const props = defineProps({
  item: {
    type: Object,
    required: false,
    default: null
  }
});

const isEditing = ref(false);
const editUrl = ref('');
const editImageUrl = ref('');
const editTags = ref('');

// itemが変わったらフォームをリセット
watch(() => props.item, (newItem) => {
  isEditing.value = false;
  if (newItem) {
    editUrl.value = newItem.url || '';
    editImageUrl.value = newItem.imageUrl || '';
    editTags.value = newItem.tags ? newItem.tags.join(', ') : '';
  }
}, { immediate: true });

const startEditing = () => {
  if (props.item) {
    editUrl.value = props.item.url || '';
    editImageUrl.value = props.item.imageUrl || '';
    editTags.value = props.item.tags ? props.item.tags.join(', ') : '';
    isEditing.value = true;
  }
};

const cancelEditing = () => {
  isEditing.value = false;
};

const openUrl = (url) => {
  BrowserOpenURL(url);
};

const saveChanges = async () => {
  if (!props.item) return;
  
  const tagsArray = editTags.value.split(',').map(t => t.trim()).filter(t => t);
  
  try {
    await UpdateProduct(props.item.path, editUrl.value, editImageUrl.value, tagsArray);
    // alert('保存しました'); // ユーザー体験を損なうのでアラートは削除してもいいかも
    
    // プロパティ更新
    props.item.url = editUrl.value;
    props.item.imageUrl = editImageUrl.value;
    props.item.tags = tagsArray;
    
    isEditing.value = false;
  } catch (error) {
    console.error('Failed to save:', error);
    alert('保存に失敗しました');
  }
};
</script>

<style scoped>
.file-preview {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 16px;
  box-sizing: border-box;
  background: white;
  border-left: 1px solid #e0e0e0;
  overflow-y: auto;
}

.empty-state {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #888;
  background: #f9f9f9;
  border-left: 1px solid #e0e0e0;
}

.preview-header {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.preview-icon {
  font-size: 32px;
  margin-right: 12px;
}

.preview-title h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.file-type {
  margin: 0;
  font-size: 12px;
  color: #888;
}

.preview-image-container {
  width: 100%;
  height: 200px;
  background-color: #f5f5f5;
  border-radius: 4px;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border: 1px solid #eee;
}

.preview-image {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.no-image {
  color: #ccc;
  font-size: 14px;
}

.preview-details {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-item label {
  font-size: 12px;
  font-weight: 600;
  color: #666;
}

.detail-item input {
  padding: 6px 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 13px;
}
.detail-item input:focus {
  border-color: #0078d7;
  outline: none;
}

.input-with-button {
  display: flex;
  gap: 4px;
}
.input-with-button input {
  flex: 1;
}
.input-with-button button {
  padding: 0 8px;
  border: 1px solid #ccc;
  background: #f0f0f0;
  border-radius: 4px;
  cursor: pointer;
}

.path-text {
  font-size: 11px;
  color: #999;
  word-break: break-all;
  background: #f5f5f5;
  padding: 4px;
  border-radius: 4px;
}

.link-text {
  color: #0078d7;
  text-decoration: none;
  font-size: 13px;
}
.link-text:hover {
  text-decoration: underline;
}

.tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.tag {
  background: #f0f0f0;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
  color: #333;
}

.actions {
  margin-top: 10px;
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.edit-btn {
  background-color: transparent;
  color: #0078d7;
  border: 1px solid #0078d7;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}
.edit-btn:hover {
  background-color: #f0f8ff;
}

.save-btn {
  background-color: #0078d7;
  color: white;
  border: none;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 600;
  font-size: 12px;
}
.save-btn:hover {
  background-color: #0063b1;
}

.cancel-btn {
  background-color: transparent;
  color: #666;
  border: 1px solid #ccc;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}
.cancel-btn:hover {
  background-color: #f5f5f5;
}
</style>
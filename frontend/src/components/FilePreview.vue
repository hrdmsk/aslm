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
  padding: 24px;
  box-sizing: border-box;
  background: white;
  border-left: 1px solid #e2e8f0;
  overflow-y: auto;
  width: 320px; /* Fixed width for better layout */
}

.empty-state {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #94a3b8;
  background: #f8fafc;
  border-left: 1px solid #e2e8f0;
  font-weight: 500;
}

.preview-header {
  display: flex;
  align-items: center;
  margin-bottom: 24px;
}

.preview-icon {
  font-size: 32px;
  margin-right: 16px;
  background-color: #f1f5f9;
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
}

.preview-title h2 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #1e293b;
  line-height: 1.4;
}

.file-type {
  margin: 2px 0 0;
  font-size: 12px;
  color: #64748b;
}

.preview-image-container {
  width: 100%;
  aspect-ratio: 16/9;
  background-color: #f8fafc;
  border-radius: 12px;
  margin-bottom: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border: 1px solid #e2e8f0;
}

.preview-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.no-image {
  color: #cbd5e1;
  font-size: 14px;
  font-weight: 500;
}

.preview-details {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.detail-item label {
  font-size: 11px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.detail-item input {
  padding: 8px 12px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  font-size: 14px;
  color: #1e293b;
  transition: all 0.2s;
  background-color: #f8fafc;
}
.detail-item input:focus {
  border-color: #6366f1;
  background-color: #ffffff;
  outline: none;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.input-with-button {
  display: flex;
  gap: 8px;
}
.input-with-button input {
  flex: 1;
}
.input-with-button button {
  padding: 0 10px;
  border: 1px solid #e2e8f0;
  background: #ffffff;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}
.input-with-button button:hover {
  background: #f1f5f9;
  border-color: #cbd5e1;
}

.path-text {
  font-size: 12px;
  color: #64748b;
  word-break: break-all;
  background: #f1f5f9;
  padding: 8px;
  border-radius: 8px;
  font-family: monospace;
}

.link-text {
  color: #6366f1;
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
}
.link-text:hover {
  text-decoration: underline;
}

.tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.tag {
  background: #f1f5f9;
  padding: 4px 10px;
  border-radius: 9999px;
  font-size: 12px;
  color: #475569;
  font-weight: 500;
}

.actions {
  margin-top: 12px;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.edit-btn {
  background-color: transparent;
  color: #6366f1;
  border: 1px solid #e0e7ff;
  padding: 8px 16px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  transition: all 0.2s;
}
.edit-btn:hover {
  background-color: #eef2ff;
  border-color: #c7d2fe;
}

.save-btn {
  background-color: #6366f1;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  font-size: 13px;
  transition: all 0.2s;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}
.save-btn:hover {
  background-color: #4f46e5;
  box-shadow: 0 4px 6px -1px rgba(99, 102, 241, 0.2);
}

.cancel-btn {
  background-color: transparent;
  color: #64748b;
  border: 1px solid #e2e8f0;
  padding: 8px 16px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  transition: all 0.2s;
}
.cancel-btn:hover {
  background-color: #f8fafc;
  color: #1e293b;
  border-color: #cbd5e1;
}
</style>
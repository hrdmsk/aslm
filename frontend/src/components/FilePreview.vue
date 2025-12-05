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

    <!-- メインプレビュー画像（親 or ファイル） -->
    <div class="preview-image-container">
      <img v-if="mainImage" :src="mainImage" alt="Preview" class="preview-image" />
      <div v-else class="no-image">No Image</div>
    </div>

    <!-- 閲覧モード -->
    <div class="preview-details" v-if="!isEditing">
      <!-- 商品フォルダのコンテキスト内では表示を最小化：商品画像・URL・編集のみ -->
      <div v-if="isInProductContext">
        <div v-if="productInfo" class="parent-info-section compact">
          <div class="section-caption">📂 商品情報</div>
          <div class="row">
            <div class="left">
              <div v-if="productInfo.url" class="link-inline">
                <a :href="productInfo.url" target="_blank" class="link-text">🔗 商品ページ</a>
              </div>
              <div v-if="productInfo.shopName" class="shop-inline">
                <small class="hint-text">店舗: {{ productInfo.shopName }}</small>
              </div>
            </div>
            <div class="right" v-if="productInfo.imageUrl">
              <img :src="productInfo.imageUrl" alt="Parent" class="small-thumb" />
            </div>
          </div>
        </div>

        <div class="actions">
          <button class="edit-btn" @click="startEditing">編集 ✎</button>
        </div>
      </div>

      <!-- 通常（商品コンテキスト外）は従来のフル表示 -->
      <div v-else>
        <!-- 親フォルダ情報（コンパクト） -->
        <div v-if="productInfo" class="parent-info-section compact">
          <div class="section-caption">📂 親フォルダ</div>
          <div class="row">
            <div class="left">
              <div v-if="productInfo.url" class="link-inline">
                <a :href="productInfo.url" target="_blank" class="link-text">🔗 商品ページ</a>
              </div>
              <div v-if="productInfo.tags && productInfo.tags.length" class="tags-inline">
                <span v-for="tag in productInfo.tags" :key="tag" class="tag small">{{ tag }}</span>
              </div>
            </div>
            <div class="right" v-if="productInfo.imageUrl">
              <img :src="productInfo.imageUrl" alt="Parent" class="small-thumb" />
            </div>
          </div>
        </div>

        <!-- ファイル情報（簡潔） -->
        <div v-if="hasFileInfo" class="file-info-section compact">
          <div class="section-caption">📄 ファイル</div>
          <div class="row">
            <div class="left">
              <div v-if="item.url" class="link-inline">
                <a :href="item.url" target="_blank" class="link-text">🔗 商品ページ（ファイル）</a>
              </div>
              <div v-if="item.tags && item.tags.length" class="tags-inline">
                <span v-for="tag in item.tags" :key="tag" class="tag small">{{ tag }}</span>
              </div>
            </div>
            <div class="right" v-if="fileThumb">
              <img :src="fileThumb" alt="File" class="small-thumb" />
            </div>
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

      <!-- Auto-fetch buttons -->
      <div class="auto-fetch-section">
        <button 
          class="auto-fetch-btn gemini-btn" 
          @click="searchProductUrl" 
          :disabled="isFetching"
        >
          <span v-if="!isFetching">✨ Geminiで調べる</span>
          <span v-else>🔄 解析中...</span>
        </button>
        
        <p v-if="fetchError" class="error-message">{{ fetchError }}</p>
        <p class="hint-text">※AIがフォルダ名から最適なBooth商品を検索します</p>
      </div>

      <div class="detail-item">
        <label>Thumbnail URL</label>
        <input type="text" v-model="editImageUrl" placeholder="https://..." />
      </div>

      <div class="detail-item">
        <label>Tags (comma separated)</label>
        <input type="text" v-model="editTags" placeholder="avatar, cute, 3d" />
      </div>

      <div class="detail-item">
        <label>Shop Name</label>
        <input type="text" v-model="editShopName" placeholder="作者/ショップ名" />
      </div>

      <div class="actions">
        <button class="cancel-btn" @click="cancelEditing">キャンセル</button>
        <button class="save-btn" @click="saveChanges" :disabled="registrationBlocked">保存</button>
      </div>
      <div v-if="registrationBlocked" class="error-message" style="margin-top:6px;">このフォルダは既に商品フォルダの子孫です。商品フォルダを作成・編集できません。</div>
    </div>
  </div>
  <div v-else class="empty-state">
    アイテムを選択してください
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue';
import { UpdateProduct, FetchBoothInfoWithGemini, FetchBoothImageFromURL, GetParentProduct, GetProductByPath } from '../../wailsjs/go/main/App';
import { BrowserOpenURL } from '../../wailsjs/runtime';
import { useFileSystemStore } from '../stores/fileSystem';

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
const editShopName = ref('');
const isFetching = ref(false);
const fetchError = ref('');

const fileSystem = useFileSystemStore();

// ローカルで選択アイテムに対する最寄りの親フォルダ情報を保持する
const nearestParent = ref(null);

// パスが祖先（ancestor が target の祖先ディレクトリか）か判定するユーティリティ
function isAncestor(ancestor, target) {
  if (!ancestor || !target) return false;
  const a = ancestor.replace(/\//g, '\\').toLowerCase();
  const t = target.replace(/\//g, '\\').toLowerCase();
  if (a === t) return true;
  // Ensure trailing separator for proper startsWith checks
  const sep = a.endsWith('\\') ? '' : '\\';
  return t.startsWith(a + sep);
}

// 商品フォルダ情報を正規化して返す（store の product context を優先）
function normalizeInfo(info) {
  if (!info) return null;
  // Support different casing coming from Go bindings
  return {
    path: info.path || info.Path || info.Pathname || null,
    name: info.name || info.Name || null,
    url: info.url || info.Url || info.URL || null,
    imageUrl: info.imageUrl || info.ImageUrl || info.ImageURL || info.image_url || null,
    shopName: info.shop || info.Shop || info.ShopName || info.shopName || null,
    tags: info.tags || info.Tags || []
  };
}

const productInfo = computed(() => {
  // 1) product context from store
  if (fileSystem.parentProductInfo) {
    const n = normalizeInfo(fileSystem.parentProductInfo);
    if (n) return n;
  }
  // 2) nearest parent fetched for selected item
  if (nearestParent.value) {
    const n = normalizeInfo(nearestParent.value);
    if (n) return n;
  }
  return null;
});

// 選択アイテムが変わったら、そのパスに対する最近接の親プロダクト情報を取得する
watch(() => props.item, async (newItem) => {
  nearestParent.value = null;
  if (!newItem || !newItem.path) return;
  try {
    // まず選択アイテム自身に紐づくプロダクト情報を取得
    let info = null;
    try {
      info = await GetProductByPath(newItem.path);
    } catch (e) {
      // GetProductByPath が存在しない、またはエラーの場合は無視して GetParentProduct を試す
      info = null;
    }

    // アイテム自身に情報が無ければ親を探す
    if (!info) {
      try {
        info = await GetParentProduct(newItem.path);
      } catch (e) {
        info = null;
      }
    }

    nearestParent.value = info || null;
  } catch (err) {
    console.error('Failed to fetch nearest parent/product for item:', err);
    nearestParent.value = null;
  }
}, { immediate: true });

// ファイル自体が情報を持っているか
const hasFileInfo = computed(() => {
  if (!props.item) return false;
  return !!(props.item.url || props.item.imageUrl || (props.item.tags && props.item.tags.length > 0));
});

// 選択アイテムが商品フォルダのコンテキスト内にいるか（商品フォルダ自身またはその子孫）
const isInProductContext = computed(() => {
  if (!productInfo.value) return false;
  if (!props.item || !props.item.path) return false;
  const p = productInfo.value.path || productInfo.value.Path || '';
  if (!p) return false;
  return isAncestor(p, props.item.path) || (p === props.item.path);
});

// 表示するメイン画像（商品フォルダ情報優先）
const mainImage = computed(() => {
  if (productInfo.value && productInfo.value.imageUrl) return productInfo.value.imageUrl;
  if (props.item && (props.item.imageUrl || props.item.ImageUrl)) return props.item.imageUrl || props.item.ImageUrl;
  return '';
});

// ファイル欄で表示するサムネイル（ファイルの imageUrl が無ければ商品フォルダの imageUrl を使う）
const fileThumb = computed(() => {
  if (!props.item) return '';
  if (props.item.imageUrl) return props.item.imageUrl;
  if (props.item.ImageUrl) return props.item.ImageUrl;
  if (productInfo.value && productInfo.value.imageUrl) return productInfo.value.imageUrl;
  return '';
});

// itemが変わったらフォームをリセット
watch(() => props.item, (newItem) => {
  isEditing.value = false;
  fetchError.value = '';
  if (newItem) {
    editUrl.value = newItem.url || '';
    editImageUrl.value = newItem.imageUrl || '';
    editTags.value = newItem.tags ? newItem.tags.join(', ') : '';
    editShopName.value = newItem.shop || newItem.Shop || '';
  }
}, { immediate: true });

const startEditing = () => {
  if (props.item) {
    // 編集時は、アイテムの情報が無ければ商品コンテキストの情報を初期値に使う
    editUrl.value = props.item.url || (productInfo.value ? productInfo.value.url : '') || '';
    editImageUrl.value = props.item.imageUrl || (productInfo.value ? productInfo.value.imageUrl : '') || '';
    editTags.value = props.item.tags ? props.item.tags.join(', ') : '';
    editShopName.value = props.item.shop || (productInfo.value ? productInfo.value.shopName : '') || '';
    isEditing.value = true;
    fetchError.value = '';
  }
};

const cancelEditing = () => {
  isEditing.value = false;
  fetchError.value = '';
};

const openUrl = (url) => {
  BrowserOpenURL(url);
};

const searchProductUrl = async () => {
  if (!props.item) return;
  
  fetchError.value = '';
  isFetching.value = true;
  
  try {
    // Use Gemini API to search
    const info = await FetchBoothInfoWithGemini(props.item.name);
    editUrl.value = info.productUrl || '';
    editImageUrl.value = info.imageUrl || '';
    editShopName.value = info.shopName || '';

    // Update the item properties immediately for preview
    props.item.url = info.productUrl || '';
    props.item.imageUrl = info.imageUrl || '';
    props.item.shop = info.shopName || '';
  } catch (error) {
    console.error('Failed to search Booth:', error);
    if (error.message && error.message.includes('API key not configured')) {
      fetchError.value = 'Gemini API Keyが設定されていません。設定画面で設定してください。';
    } else {
      fetchError.value = '商品URLの検索に失敗しました';
    }
  } finally {
    isFetching.value = false;
  }
};

const saveChanges = async () => {
  if (!props.item) return;
  
  const tagsArray = editTags.value.split(',').map(t => t.trim()).filter(t => t);
  
  try {
    await UpdateProduct(props.item.path, editUrl.value, editImageUrl.value, editShopName.value, tagsArray);

    // プロパティ更新
    props.item.url = editUrl.value;
    props.item.imageUrl = editImageUrl.value;
    props.item.tags = tagsArray;
    props.item.shop = editShopName.value;
    
    isEditing.value = false;
    fetchError.value = '';
    // 更新後、DB側の最新の product info を取得してストアとローカル nearestParent を更新
    try {
      let updated = null;
      try {
        updated = await GetProductByPath(props.item.path);
      } catch (e) {
        updated = null;
      }
      if (!updated) {
        try {
          updated = await GetParentProduct(props.item.path);
        } catch (e) {
          updated = null;
        }
      }
      if (updated) {
        nearestParent.value = updated;
        // update store product context when appropriate
        if (fileSystem.parentProductInfo) {
          const storePath = fileSystem.parentProductInfo.Path || fileSystem.parentProductInfo.path || '';
          const updPath = updated.Path || updated.path || '';
          if (storePath === updPath || isAncestor(updPath, storePath)) {
            fileSystem.parentProductInfo = updated;
          }
        }
      }
    } catch (e) {
      console.warn('Failed to refresh product info after save:', e);
    }
  } catch (error) {
    console.error('Failed to save:', error);
    alert('保存に失敗しました');
  }
};

// 登録（保存）をブロックすべきか判定（商品フォルダの子孫ならブロック）
const registrationBlocked = computed(() => {
  if (!props.item) return false;
  if (props.item.type !== 'folder') return false; // ファイルは問題ない

  // parent product from store (product context)
  const p = fileSystem.parentProductInfo;
  const parentPath = p ? (p.path || p.Path || null) : null;

  if (!parentPath) return false;

  // If the parentPath equals the current folder, allow editing (it's the product folder itself)
  const itemPath = props.item.path || props.item.Path || '';
  if (!itemPath) return false;

  if (parentPath === itemPath) return false;

  // ancestor check
  const a = parentPath.replace(/\//g, '\\').toLowerCase();
  const t = itemPath.replace(/\//g, '\\').toLowerCase();
  const sep = a.endsWith('\\') ? '' : '\\';
  return t.startsWith(a + sep);
});
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
  width: 320px;
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
  height: 120px;
  background-color: #f8fafc;
  border-radius: 8px;
  margin-bottom: 12px;
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

.auto-fetch-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.auto-fetch-btn {
  color: white;
  border: none;
  padding: 12px 16px;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.gemini-btn {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  box-shadow: 0 2px 4px rgba(79, 172, 254, 0.2);
}
.gemini-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(79, 172, 254, 0.3);
}

.auto-fetch-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.hint-text {
  font-size: 11px;
  color: #94a3b8;
  margin: 0;
  font-style: italic;
}

.error-message {
  color: #ef4444;
  font-size: 12px;
  margin: 0;
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

.section-caption {
  font-size: 12px;
  color: #334155;
  font-weight: 600;
  margin-bottom: 6px;
}

.parent-info-section.compact,
.file-info-section.compact {
  padding: 8px;
  background: #fff;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  border: 1px solid #eef2f7;
}

.row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.left { flex: 1; }
.right { margin-left: 8px; }
.small-thumb { width: 72px; height: 48px; object-fit: cover; border-radius: 6px; border: 1px solid #e2e8f0; }
.tag.small { padding: 2px 6px; font-size: 11px; }
.link-inline { margin-bottom: 6px; }

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
<template>
  <div class="address-bar-container">
    <!-- ナビゲーションボタン群 -->
    <div class="nav-buttons">
      <button class="nav-btn" @click="goBack" title="戻る">←</button>
      <button class="nav-btn" @click="goForward" title="進む" disabled>→</button>
      <button class="nav-btn" @click="goUp" title="上の階層へ">↑</button>
    </div>

    <!-- アドレス入力エリア -->
    <div class="address-input-wrapper">
      <div class="icon-folder">📁</div>
      <!-- キーボード入力でパス移動できるようにenterイベントを設定 -->
      <input 
        type="text" 
        class="address-input" 
        :value="store.currentPath" 
        @keydown.enter="handleManualInput"
      />
    </div>

    <!-- 検索ボックス -->
    <div class="search-box">
      <input type="text" placeholder="ASLMの検索" />
      <span class="search-icon">🔍</span>
    </div>
  </div>
</template>

<script setup>
import { useFileSystemStore } from '../stores/fileSystem';

const store = useFileSystemStore();

// ※ 履歴管理（戻る/進む）は本来Routerや履歴スタックが必要ですが、ここでは簡易実装です
const goBack = () => {
  console.log('History back (Not implemented yet)');
};

const goForward = () => {
  console.log('History forward (Not implemented yet)');
};

// 親ディレクトリへ移動
const goUp = () => {
  const current = store.currentPath;
  // シンプルな文字列操作で親パスを計算 (Windows/Unixパス区切りに対応が必要)
  const separator = current.includes('/') ? '/' : '\\';
  const parts = current.split(separator);
  
  if (parts.length > 1) {
    parts.pop(); // 末尾を削除
    // ルート直下(C:/)の場合の空文字対策などは適宜必要
    const newPath = parts.join(separator) || separator; 
    store.changeDirectory(newPath);
  }
};

const handleManualInput = (e) => {
  store.changeDirectory(e.target.value);
};
</script>

<style scoped>
.address-bar-container {
  display: flex;
  align-items: center;
  width: 100%;
  gap: 12px;
}

.nav-buttons {
  display: flex;
  gap: 4px;
}

.nav-btn {
  background: transparent;
  border: none;
  font-size: 16px;
  color: #555;
  padding: 4px 8px;
  border-radius: 4px;
  cursor: pointer;
}
.nav-btn:hover { background-color: #e5e5e5; }
.nav-btn:disabled { color: #ccc; cursor: default; }

.address-input-wrapper {
  flex: 1;
  display: flex;
  align-items: center;
  background-color: white;
  border: 1px solid #d9d9d9;
  border-radius: 2px; /* エクスプローラーらしい角ばったデザイン */
  padding: 2px 6px;
  height: 28px;
  box-shadow: inset 0 1px 2px rgba(0,0,0,0.05);
}

.address-input-wrapper:focus-within {
  border-color: #0078d7; /* フォーカス時の青枠 */
}

.icon-folder {
  margin-right: 6px;
  font-size: 14px;
  color: #ffc107; /* フォルダ色 */
}

.address-input {
  border: none;
  width: 100%;
  outline: none;
  font-size: 13px;
  color: #333;
}

.search-box {
  width: 200px;
  background-color: white;
  border: 1px solid #d9d9d9;
  border-radius: 2px;
  height: 28px;
  display: flex;
  align-items: center;
  padding: 0 8px;
}
.search-box input {
  border: none;
  outline: none;
  width: 100%;
  font-size: 12px;
}
.search-icon {
  font-size: 12px;
  color: #888;
}
</style>
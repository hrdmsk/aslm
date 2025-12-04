<template>
  <div class="address-bar-container">
    <!-- ナビゲーションボタン群 -->
    <div class="nav-buttons">
      <button 
        class="nav-btn" 
        @click="store.goBack()" 
        title="戻る" 
        :disabled="store.historyIndex <= 0"
      >←</button>
      <button 
        class="nav-btn" 
        @click="store.goForward()" 
        title="進む" 
        :disabled="store.historyIndex >= store.historyStack.length - 1"
      >→</button>
      <button class="nav-btn" @click="store.goUp()" title="上の階層へ">↑</button>
      <button class="nav-btn" @click="handleRefresh()" title="更新">🔄</button>
      <button class="nav-btn" @click="handleHome()" title="ホーム">🏠</button>
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

const handleManualInput = (e) => {
  store.changeDirectory(e.target.value);
};

const handleRefresh = () => {
  store.changeDirectory(store.currentPath, false);
};

const handleHome = () => {
  store.changeDirectory(store.homePath);
};
</script>

<style scoped>
.address-bar-container {
  display: flex;
  align-items: center;
  width: 100%;
  gap: 16px;
}

.nav-buttons {
  display: flex;
  gap: 8px;
}

.nav-btn {
  background: transparent;
  border: none;
  font-size: 18px;
  color: #64748b;
  padding: 8px;
  border-radius: 50%;
  cursor: pointer;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}
.nav-btn:hover:not(:disabled) { 
  background-color: #f1f5f9; 
  color: #1e293b;
}
.nav-btn:disabled { 
  color: #cbd5e1; 
  cursor: default; 
}

.address-input-wrapper {
  flex: 1;
  display: flex;
  align-items: center;
  background-color: #f1f5f9;
  border: 1px solid transparent;
  border-radius: 9999px; /* Pill shape */
  padding: 8px 16px;
  height: 40px;
  transition: all 0.2s;
}

.address-input-wrapper:focus-within {
  background-color: #ffffff;
  border-color: #6366f1;
  box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.1);
}

.icon-folder {
  margin-right: 10px;
  font-size: 16px;
  color: #64748b;
}

.address-input {
  border: none;
  width: 100%;
  outline: none;
  font-size: 14px;
  color: #1e293b;
  background: transparent;
}

.search-box {
  width: 240px;
  background-color: #f1f5f9;
  border: 1px solid transparent;
  border-radius: 9999px;
  height: 40px;
  display: flex;
  align-items: center;
  padding: 0 16px;
  transition: all 0.2s;
}
.search-box:focus-within {
  background-color: #ffffff;
  border-color: #6366f1;
  box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.1);
}

.search-box input {
  border: none;
  outline: none;
  width: 100%;
  font-size: 14px;
  background: transparent;
  color: #1e293b;
}
.search-icon {
  font-size: 14px;
  color: #64748b;
  margin-left: 8px;
}
</style>
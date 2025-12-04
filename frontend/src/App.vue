<template>
  <div class="app-layout">
    <header class="app-header">
      <AddressBar />
      <div class="settings-button">
        <button @click="showSettings = !showSettings" title="設定">⚙️</button>
      </div>
    </header>

    <div class="app-body">
      <aside class="sidebar">
        <div class="sidebar-placeholder">Quick Access (Tree View)</div>
      </aside>

      <main class="content">
        <SettingsScreen v-if="showSettings" @close="showSettings = false" />
        <FileList v-else />
      </main>
    </div>
    
    <footer class="app-footer">
      {{ store.fileList.length }} 個の項目
    </footer>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import FileList from './components/FileList.vue';
import SettingsScreen from './components/SettingsScreen.vue';
import AddressBar from './components/AddressBar.vue';
import { useFileSystemStore } from './stores/fileSystem';

const store = useFileSystemStore();
const showSettings = ref(false);

// 初期ロード - ホームディレクトリから開始
store.changeDirectory(store.homePath);
</script>

<style>
/* リセットCSS等は省略。全体レイアウト */
html, body, #app { height: 100%; margin: 0; overflow: hidden; }

.app-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}

.app-header {
  height: 64px;
  background-color: #ffffff;
  border-bottom: 1px solid #e2e8f0;
  display: flex;
  align-items: center;
  padding: 0 24px;
  gap: 16px;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  z-index: 10;
}

.app-body {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.sidebar {
  width: 240px;
  background-color: #ffffff;
  border-right: 1px solid #e2e8f0;
  padding: 20px;
  display: flex;
  flex-direction: column;
}

.sidebar-placeholder {
  color: #64748b;
  font-size: 14px;
  font-weight: 500;
}

.content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: #f8fafc;
  position: relative;
}

.app-footer {
  height: 32px;
  background-color: #ffffff;
  border-top: 1px solid #e2e8f0;
  display: flex;
  align-items: center;
  padding: 0 16px;
  font-size: 12px;
  color: #64748b;
}

.settings-button button {
  background: transparent;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: #64748b;
  padding: 8px;
  border-radius: 50%;
  transition: background-color 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.settings-button button:hover {
  background-color: #f1f5f9;
  color: #1e293b;
}
</style>
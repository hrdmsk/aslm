<template>
  <div class="app-layout">
    <header class="app-header">
      <div class="nav-buttons">
        <button>←</button> <button>→</button> <button>↑</button>
      </div>
      <div class="address-bar">
        <span>📁</span>
        <input type="text" :value="store.currentPath" readonly />
      </div>
      <div class="search-bar">
        <input type="text" placeholder="ASLMの検索" />
      </div>
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
import { useFileSystemStore } from './stores/fileSystem';

const store = useFileSystemStore();
const showSettings = ref(false);

// 初期ロード
store.changeDirectory('C:/UnityAssets');
</script>

<style>
/* リセットCSS等は省略。全体レイアウト */
html, body, #app { height: 100%; margin: 0; overflow: hidden; }

.app-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  color: #333;
}

.app-header {
  height: 40px;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 10px;
  border-bottom: 1px solid #d9d9d9;
  background: #f3f3f3;
}

.address-bar {
  flex: 1;
  display: flex;
  align-items: center;
  background: white;
  border: 1px solid #ccc;
  padding: 4px 8px;
  border-radius: 2px;
}
.address-bar input { border: none; width: 100%; outline: none; font-size: 13px; }

.app-body {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.sidebar {
  width: 250px;
  border-right: 1px solid #d9d9d9;
  background: #f0f0f0;
  overflow-y: auto;
  padding: 10px;
}

.content {
  flex: 1;
  overflow: hidden;
}

.app-footer {
  height: 24px;
  background: #f3f3f3;
  border-top: 1px solid #d9d9d9;
  display: flex;
  align-items: center;
  padding: 0 10px;
  font-size: 11px;
  color: #666;
}
</style>
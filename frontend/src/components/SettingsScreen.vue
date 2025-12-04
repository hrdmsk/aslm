<template>
  <div class="settings-container">
    <div class="settings-card">
      <h2>設定</h2>
      <div class="setting-item">
        <label for="home-path">ホームディレクトリ</label>
        <div class="input-group">
          <input 
            id="home-path" 
            type="text" 
            v-model="localHomePath" 
            placeholder="C:/UnityAssets"
          />
        </div>
      </div>
      <div class="actions" style="display: flex; justify-content: flex-end; gap: 12px;">
        <button @click="$emit('close')" class="save-btn" style="background-color: transparent; color: #64748b; border: 1px solid #e2e8f0;">キャンセル</button>
        <button @click="saveSettings" class="save-btn">保存</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useFileSystemStore } from '../stores/fileSystem';

const store = useFileSystemStore();
const localHomePath = ref(store.homePath);
const emit = defineEmits(['close']);

const saveSettings = () => {
  store.setHomePath(localHomePath.value);
  emit('close');
};
</script>

<style scoped>
.settings-container {
  padding: 40px;
  max-width: 600px;
  margin: 0 auto;
  font-family: 'Inter', 'Segoe UI', sans-serif;
}

.settings-card {
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  padding: 32px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

h2 {
  margin-top: 0;
  margin-bottom: 24px;
  font-size: 24px;
  font-weight: 600;
  color: #1e293b;
}

.setting-item {
  margin-bottom: 24px;
}

.setting-item label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #475569;
  font-size: 14px;
}

.input-group {
  display: flex;
  gap: 12px;
}

.input-group input {
  flex: 1;
  padding: 10px 16px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  font-size: 14px;
  color: #1e293b;
  transition: all 0.2s;
  background-color: #f8fafc;
}

.input-group input:focus {
  border-color: #6366f1;
  background-color: #ffffff;
  outline: none;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.save-btn {
  padding: 10px 20px;
  background-color: #6366f1;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  font-size: 14px;
  transition: all 0.2s;
  white-space: nowrap;
}

.save-btn:hover {
  background-color: #4f46e5;
  box-shadow: 0 4px 6px -1px rgba(99, 102, 241, 0.2);
}

.save-btn:active {
  transform: translateY(0);
}
</style>

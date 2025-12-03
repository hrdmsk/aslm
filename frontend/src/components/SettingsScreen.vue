<template>
  <div class="settings-screen">
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
    <div class="actions">
      <button @click="saveSettings" class="save-btn">保存</button>
      <button @click="$emit('close')" class="cancel-btn">キャンセル</button>
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
.settings-screen {
  padding: 20px;
  background: #fff;
  height: 100%;
  box-sizing: border-box;
}

h2 {
  margin-top: 0;
  margin-bottom: 20px;
  font-size: 20px;
  color: #333;
}

.setting-item {
  margin-bottom: 20px;
}

.setting-item label {
  display: block;
  margin-bottom: 8px;
  font-weight: bold;
  color: #555;
}

.input-group input {
  width: 100%;
  max-width: 400px;
  padding: 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 14px;
}

.actions {
  display: flex;
  gap: 10px;
}

button {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.save-btn {
  background-color: #007bff;
  color: white;
}

.save-btn:hover {
  background-color: #0056b3;
}

.cancel-btn {
  background-color: #f0f0f0;
  color: #333;
}

.cancel-btn:hover {
  background-color: #e0e0e0;
}
</style>

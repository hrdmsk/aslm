import { defineStore } from 'pinia';
import { ref } from 'vue';
import { ListFiles } from '../../wailsjs/go/main/App';

export const useFileSystemStore = defineStore('fileSystem', () => {
  const currentPath = ref('D:/VRChatAssetPack');
  const fileList = ref([]);
  const directoryTree = ref([]);
  const parentProductInfo = ref(null);

  // 履歴管理用
  const historyStack = ref([currentPath.value]);
  const historyIndex = ref(0);

  async function changeDirectory(path, addToHistory = true) {
    try {
      const files = await ListFiles(path);
      fileList.value = files || [];
      currentPath.value = path;

      // 履歴に追加
      if (addToHistory) {
        // 現在のインデックスより後ろの履歴は削除
        historyStack.value = historyStack.value.slice(0, historyIndex.value + 1);
        historyStack.value.push(path);
        historyIndex.value++;
      }
    } catch (err) {
      console.error('Error changing directory:', err);
    }
  }

  function goBack() {
    if (historyIndex.value > 0) {
      historyIndex.value--;
      const path = historyStack.value[historyIndex.value];
      changeDirectory(path, false);
    }
  }

  function goForward() {
    if (historyIndex.value < historyStack.value.length - 1) {
      historyIndex.value++;
      const path = historyStack.value[historyIndex.value];
      changeDirectory(path, false);
    }
  }

  function goUp() {
    const current = currentPath.value;
    const separator = current.includes('/') ? '/' : '\\';
    const parts = current.split(separator);

    // 末尾が空文字（セパレータで終わる場合）の対策
    if (parts[parts.length - 1] === '') parts.pop();

    if (parts.length > 1) {
      parts.pop(); // 末尾を削除
      // ドライブ直下(C:)などの場合、セパレータをつける必要があるかも
      let newPath = parts.join(separator);
      if (!newPath.includes(separator) && !newPath.includes('\\')) {
        newPath += separator; // C: -> C:/
      }
      changeDirectory(newPath);
    }
  }

  // ホームディレクトリの設定
  const homePath = ref('D:/VRChatAssetPack');

  function setHomePath(path) {
    homePath.value = path;
    // 必要であればここで永続化処理などを呼び出す
    console.log(`Home path set to: ${path}`);
    // ホームディレクトリ変更時にそのディレクトリを表示するかは要件次第だが、
    // ここでは設定変更だけにしておく
  }

  return {
    currentPath,
    fileList,
    directoryTree,
    changeDirectory,
    homePath,
    setHomePath,
    historyStack,
    historyIndex,
    goBack,
    goForward,
    goUp,
    parentProductInfo
  };
});
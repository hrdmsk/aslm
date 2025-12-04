import { defineStore } from 'pinia';
import { ref } from 'vue';
import { ListFiles } from '../../wailsjs/go/main/App';

export const useFileSystemStore = defineStore('fileSystem', () => {
  // 現在のパス
  const currentPath = ref('D:/VRChatAssetPack');

  // 現在のディレクトリ内のファイルリスト
  const fileList = ref([]);

  // 左側のツリー構造データ
  const directoryTree = ref([]);

  // 履歴管理
  const historyStack = ref([]);
  const historyIndex = ref(-1);

  // アクション: パスを変更してGoバックエンドにリクエストを送る
  async function changeDirectory(path, addToHistory = true) {
    // 同じパスなら何もしない（ただし履歴移動の場合は再読み込みしたいかも？一旦スキップ）
    if (currentPath.value === path && addToHistory) return;

    console.log(`Fetching files for: ${path}`);
    try {
      const files = await ListFiles(path);
      fileList.value = files || [];
      currentPath.value = path;

      if (addToHistory) {
        // 履歴の現在位置より先を削除して新しいパスを追加
        if (historyIndex.value < historyStack.value.length - 1) {
          historyStack.value = historyStack.value.slice(0, historyIndex.value + 1);
        }
        historyStack.value.push(path);
        historyIndex.value++;
      }
    } catch (error) {
      console.error('Failed to list files:', error);
      fileList.value = [];
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
    goUp
  };
});
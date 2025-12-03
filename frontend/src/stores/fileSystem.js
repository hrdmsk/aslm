import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useFileSystemStore = defineStore('fileSystem', () => {
  // 現在のパス
  const currentPath = ref('C:/UnityAssets');

  // 現在のディレクトリ内のファイルリスト（Goから取得するデータ）
  // ユーザー要望に合わせてフォルダにもURLとSourceを付与
  const fileList = ref([
    {
      name: 'Avatar_Model_v1.unitypackage',
      type: 'file',
      source: 'booth',
      url: 'https://booth.pm/ja/items/123456'
    },
    {
      name: 'Shader_Pack_Vol1',
      type: 'folder',
      source: 'booth', // フォルダにもソース情報を付与
      url: 'https://booth.pm/ja', // ユーザー指定のサンプルURL
      imageUrl: 'https://go.dev/images/gophers/motorcycle.svg' // フォルダのプレビュー画像をURLで紐付け
    },
    {
      name: 'Prop_Set_v2.unitypackage',
      type: 'file',
      source: 'gumroad',
      url: 'https://gumroad.com/l/sample'
    },
    {
      name: 'Environment_Assets',
      type: 'folder',
      source: null,
      url: ''
    },
  ]);

  // 左側のツリー構造データ（ダミー）
  const directoryTree = ref([
    {
      name: 'UnityAssets',
      path: 'C:/UnityAssets',
      isExpanded: true,
      children: [
        { name: 'Shader_Pack_Vol1', path: 'C:/UnityAssets/Shader_Pack_Vol1', isExpanded: false, children: [] },
        { name: 'Environment_Assets', path: 'C:/UnityAssets/Environment_Assets', isExpanded: false, children: [] },
      ]
    }
  ]);

  // アクション: パスを変更してGoバックエンドにリクエストを送る（擬似コード）
  async function changeDirectory(path) {
    currentPath.value = path;
    console.log(`Fetching files for: ${path}`);
    // 実際にはここでAPIを叩いて fileList を更新する
    // 今回はモックなのでリストの内容は固定のままにしています
  }

  // ホームディレクトリの設定
  const homePath = ref('C:/UnityAssets');

  function setHomePath(path) {
    homePath.value = path;
    // 必要であればここで永続化処理などを呼び出す
    console.log(`Home path set to: ${path}`);
  }

  return { currentPath, fileList, directoryTree, changeDirectory, homePath, setHomePath };
});
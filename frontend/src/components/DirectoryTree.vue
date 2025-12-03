<template>
  <ul class="tree-root">
    <li v-for="node in nodes" :key="node.path">
      <!-- フォルダ行本体 -->
      <div 
        class="tree-item" 
        :class="{ 'is-active': store.currentPath === node.path }"
        @click="selectDirectory(node)"
      >
        <!-- 展開/折りたたみ矢印 (子要素がある場合のみ表示) -->
        <span 
          class="toggle-icon" 
          @click.stop="toggleExpand(node)"
          :class="{ 'is-hidden': !node.children || node.children.length === 0 }"
        >
          {{ node.isExpanded ? '▼' : '▶' }}
        </span>

        <!-- フォルダアイコン -->
        <span class="folder-icon">
           {{ node.isExpanded ? '📂' : '📁' }}
        </span>

        <!-- フォルダ名 -->
        <span class="folder-name">{{ node.name }}</span>
      </div>

      <!-- 再帰呼び出し: 展開されている場合のみ子要素を表示 -->
      <div v-if="node.isExpanded && node.children && node.children.length > 0" class="tree-children">
        <!-- 自分自身(DirectoryTree)を呼び出す -->
        <DirectoryTree :nodes="node.children" />
      </div>
    </li>
  </ul>
</template>

<script setup>
import { useFileSystemStore } from '../stores/fileSystem';

// props定義: 再帰的に渡されるノードの配列を受け取る
const props = defineProps({
  nodes: {
    type: Array,
    default: () => []
  }
});

const store = useFileSystemStore();

// ディレクトリ選択時の処理
const selectDirectory = (node) => {
  store.changeDirectory(node.path);
};

// 展開/折りたたみトグル
const toggleExpand = (node) => {
  // 注意: 本来はStoreのアクション経由で状態を変更すべきですが、
  // ここでは簡易的にオブジェクトのプロパティを直接変更しています。
  // 実装時はデータのリアクティブ性を確保してください。
  node.isExpanded = !node.isExpanded;
};
</script>

<style scoped>
.tree-root {
  list-style: none;
  padding: 0;
  margin: 0;
  font-size: 13px;
}

.tree-item {
  display: flex;
  align-items: center;
  padding: 4px 2px;
  cursor: default;
  white-space: nowrap;
  user-select: none;
}

.tree-item:hover {
  background-color: #f0f9ff; /* ホバー色 */
}

/* 選択中のディレクトリ */
.tree-item.is-active {
  background-color: #cce8ff;
  border: 1px dotted #99d1ff; /* Windowsっぽい点線枠 */
}

.toggle-icon {
  width: 16px;
  text-align: center;
  font-size: 10px;
  color: #888;
  cursor: pointer;
  margin-right: 2px;
}
.toggle-icon.is-hidden {
  opacity: 0;
  pointer-events: none;
}

.folder-icon {
  margin-right: 6px;
  color: #ffc107;
}

.folder-name {
  flex: 1;
}

/* インデント（階層構造の表現） */
.tree-children {
  padding-left: 18px; /* インデント幅 */
}
</style>
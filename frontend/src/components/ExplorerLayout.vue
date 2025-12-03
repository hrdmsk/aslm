<template>
  <div class="explorer-layout">
    <!-- 上部：アドレスバー・ツールバー -->
    <header class="layout-header">
      <slot name="header"></slot>
    </header>

    <div class="layout-body">
      <!-- 左側：ツリービュー -->
      <aside class="layout-sidebar">
        <slot name="sidebar"></slot>
      </aside>

      <!-- 右側：メインコンテンツ -->
      <main class="layout-content">
        <slot name="content"></slot>
      </main>
    </div>

    <!-- 下部：ステータスバー -->
    <footer class="layout-footer">
      <slot name="footer"></slot>
    </footer>
  </div>
</template>

<script setup>
// ロジックは持たず、レイアウト構造の定義に徹する
</script>

<style scoped>
.explorer-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100vw;
  overflow: hidden;
  background-color: #f0f0f0;
  color: #333;
  font-family: 'Segoe UI', sans-serif;
}

.layout-header {
  height: 48px; /* ツールバーを含む標準的な高さ */
  background-color: #f9f9f9;
  border-bottom: 1px solid #d9d9d9;
  display: flex;
  align-items: center;
  padding: 0 10px;
  flex-shrink: 0; /* 縮小させない */
}

.layout-body {
  display: flex;
  flex: 1; /* 残りの高さを占有 */
  overflow: hidden;
}

.layout-sidebar {
  width: 260px; /* 初期のサイドバー幅 */
  background-color: #f0f0f0; /* エクスプローラーのサイドバー色 */
  border-right: 1px solid #d9d9d9;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 5px 0;
  flex-shrink: 0;
}

.layout-content {
  flex: 1;
  background-color: #ffffff;
  overflow: hidden; /* 内部でスクロールさせるため */
  position: relative;
}

.layout-footer {
  height: 24px;
  background-color: #f0f0f0;
  border-top: 1px solid #d9d9d9;
  display: flex;
  align-items: center;
  padding: 0 12px;
  font-size: 11px;
  color: #555;
  flex-shrink: 0;
}
</style>
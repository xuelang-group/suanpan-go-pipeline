<template>
  <div class="ml-init-loading">
    <div class="ml-loading">
      <a-spin size="large" />
    </div>
  </div>
</template>

<script>
export default {
  name: 'InitialLoading',
}
</script>

<style lang="less">
.ml-init-loading {
  position: fixed;
  z-index: 9999;
  left: 0;
  top: 0;
  width: 100vw;
  height: 100vh;
  background: #f5f5f5;
}
</style>
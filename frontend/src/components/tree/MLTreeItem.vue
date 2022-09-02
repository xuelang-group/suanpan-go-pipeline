<template>
  <li class="ml-tree-item">
    <div class="ml-tree-item-inner" @click="expandHandler" @mousedown.capture="startDrag">
      <span class="ml-tree-item-placeholder" :style="{ width: `${data.depth * 10}px` }"></span>
      <span class="ml-tree-item-icon">
        <span v-if="data.isDirectory && data.expand" class="iconfont icon-arrow-down-full"></span>
        <span v-if="data.isDirectory && !data.expand" class="iconfont icon-arrow-right-full"></span>
        <span v-if="!data.isDirectory" class="iconfont icon-cube"></span>
      </span>
      {{ data.title }}
    </div>
    <ul v-show="data.isDirectory && data.expand">
      <MLTreeItem v-for="treeItem in data.children" :key="treeItem.key" :data="treeItem"></MLTreeItem>
    </ul>
  </li>
</template>

<script>
import { findComponentUpward } from '../../utils'

export default {
  name: 'MLTreeItem',
  props: {
    data: {
      type: Object
    }
  },
  data() {
    return {}
  },
  created() {
    this.treeComponent = findComponentUpward(this, 'MLTree')
  },
  methods: {
    expandHandler() {
      this.data.expand = !this.data.expand
    },
    startDrag(e) {
      if(this.data.isDirectory || (e.button !== 0)) {
        return;
      }
      this.treeComponent.startDrag(this.data, e.clientX, e.clientY)
    }
  }
}
</script>

<style>

</style>
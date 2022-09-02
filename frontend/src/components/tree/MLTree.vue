<template>
  <ul class="ml-tree">
    <MLTreeItem v-for="treeItem in treeData" :key="treeItem.key" :data="treeItem"></MLTreeItem>
  </ul>
</template>

<script>
import MLTreeItem from './MLTreeItem.vue'

export default {
  name: 'MLTree',
  components: {
    MLTreeItem
  },
  props: {
    treeData: {
      type: Array,
      default: []
    }
  },
  emits: ['start-drag'],
  watch: {
    treeData: {
      immediate: true,
      handler() {
        this.updateData()
      }
    }
  },
  data() {
    return {
      treeDataFormat: []
    }
  },
  methods: {
    updateData() {
      this.treeDataFormat = []
      this.dataFormat(this.treeData, 0, this.treeDataFormat);
    },
    dataFormat(treeData, depth, treeDataFormat) {
      for(let i = 0, len = treeData.length; i < len; i++) {
        let item = treeData[i]
        item.depth = depth
        treeDataFormat.push(item)
        item.isDirectory = false
        if(item.children) {
          item.expand = true
          item.isDirectory = true
          this.dataFormat(item.children, depth + 1, treeDataFormat)
        }
      }
    },
    startDrag(componentData, x, y) {
      this.$emit('start-drag', componentData, x, y)
    }
  }
}
</script>

<style>

</style>
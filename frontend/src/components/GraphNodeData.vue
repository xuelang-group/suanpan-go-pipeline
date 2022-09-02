<template>
  <a-modal 
    v-model:visible="modalVisible"
    :title="'数据预览-'"
    :maskClosable="false"
    :footer="null"
    :afterClose="afterCloseHandler">
  </a-modal>
</template>

<script>
export default {
  name: 'GraphNodeData',
  props: {
    visible: {
      type: Boolean,
      default: false
    }
  },
  emits: ['update:visible'],
  data() {
    return {
      modalVisible: this.visible,
    }
  },
  watch: {
    visible() {
      this.modalVisible = this.visible
    }
  },
  methods: {
    afterCloseHandler() {
      this.$emit('update:visible', false)
    },
  }
}
</script>

<style>

</style>
<template>
  <div class="ml-collapse-panel">
    <div 
      class="ml-collapse-panel-header"
      :style="{ height: `${headerHeight}px`}"
      @click="expendHandler">
      <!-- <span 
        class="ml-collapse-panel-header-icon iconfont icon-arrow-down" 
        :style="{ transform: expendPanel ? `rotate(0deg)` : 'rotate(-90deg)' }"></span> -->
      <div class="ml-collapse-panel-header-label">{{ header }}</div>
      <div 
        v-show="showResize" 
        class="ml-collapse-panel-resize" 
        @click.stop
        @mousedown.stop="resizeMousedown"></div>
    </div>
    <div 
      class="ml-collapse-panel-content"
      :style="{ height: `${contentHeight}px`, transition: resizing ? '' : 'height 0.25s' }"
      v-show="expendPanel">
      <slot></slot>
    </div>
  </div>
</template>

<script>
export default {
  name: 'MLCollapsePanel',
  props: {
    expend: {
      type: Boolean,
      default: false
    },
    header: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      expendPanel: this.expend,
      headerHeight: 0,
      contentHeight: 100,
      showResize: false,
      resizing: false
    }
  },
  watch: {
    expend(val) {
      this.expendPanel = val
    },
    expendPanel() {
      this.$parent.panelExpendChanged()
    }
  },
  created() {
    this.$parent.collapsePanels.push(this); 
  },
  beforeUnmount() {
    if(this.$parent.collapsePanels) {
      let idx = this.$parent.collapsePanels.indexOf(this)
      if(idx > -1) {
        this.$parent.collapsePanels.splice(idx, 1)
      }
    }
  },
  methods: {
    expendHandler() {
      this.expendPanel = !this.expendPanel
    },
    resizeMousedown(e) {
      this.resizing = true
      this.$parent.resizeStart(e, this)
    },
  }
}
</script>

<style lang="less">

</style>
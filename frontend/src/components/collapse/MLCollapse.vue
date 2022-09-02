<template>
  <div 
    ref="collapse" 
    class="ml-collapse">
    <slot></slot>
  </div>
</template>

<script>
import { markRaw } from 'vue'

export default {
  name: 'MLCollapse',
  data() {
    return {
      collapsePanels: markRaw([]),
      panelHeaderHeight: 44,
      resizingPanel: null,
      panelContentMinHeight: 100
    }
  },
  mounted() {
    window.addEventListener('resize', this.resizeChanged)
    this.updateHeight()
    this.initHeight()
  },
  beforeUnmount() {
    this.collapsePanels = null;
    window.removeEventListener('resize', this.resizeChanged)
  },
  methods: {
    resizeChanged() {
      this.updateHeight()
    },
    initHeight() {
      this.collapsePanels.forEach(penel => {
        penel.headerHeight = this.panelHeaderHeight
      });
      this.updateHeight()
      this.updateResize()
    },
    updateHeight() {
      if(this.collapsePanels.length < 1) {
        return
      }
      let collapseHeight = this.$refs.collapse.clientHeight
      let expendPanels = this.collapsePanels.filter(panel => panel.expendPanel)
      expendPanels.forEach(panel => {
        panel.contentHeight = (collapseHeight - 
          this.panelHeaderHeight * this.collapsePanels.length) / expendPanels.length;
      })
    },
    panelExpendChanged() {
      if(this.collapsePanels.length < 1) {
        return
      }
      this.updateHeight();
      this.updateResize();
    },
    updateResize() {
      this.collapsePanels.forEach(panel => panel.showResize = false);
      for(let i = 0; i < this.collapsePanels.length - 1; i++) {
        if(this.collapsePanels[i].expendPanel && this.collapsePanels[i+1].expendPanel) {
          this.collapsePanels[i+1].showResize = true
        }
      }
    },
    resize(panel, diff) {
      let targetIdx = this.collapsePanels.indexOf(panel)
      if(targetIdx > -1) {
        let newcontentHeight1 = this.collapsePanels[targetIdx].contentHeight - diff,
          newcontentHeight2 = this.collapsePanels[targetIdx - 1].contentHeight + diff;
        if((newcontentHeight1 < this.panelContentMinHeight) || (newcontentHeight2 < this.panelContentMinHeight)) {
          return;
        }
        this.collapsePanels[targetIdx].contentHeight = newcontentHeight1
        this.collapsePanels[targetIdx - 1].contentHeight = newcontentHeight2
      }
    },
    resizeStart(e, panel) {
      this.resizingPanel = panel
      this.resizePrev = e.clientY
      this.collapsePanels.forEach(panel => panel.resizing = true)
      window.addEventListener('mousemove', this.resizeMousemove)
      window.addEventListener('mouseup', this.resizeMouseup)
    },
    resizeMousemove(e) {
      e.preventDefault()
      if(!this.resizingPanel) {
        return;
      }
      let diff = e.clientY - this.resizePrev
      this.resizePrev = e.clientY
      this.resize(this.resizingPanel, diff)
    },
    resizeMouseup() {
      this.resizingPanel = null
      this.collapsePanels.forEach(panel => panel.resizing = false)
      window.removeEventListener('mousemove', this.resizeMousemove)
      window.removeEventListener('mouseup', this.resizeMouseup)
    }
  }
}
</script>

<style lang="less">

</style>
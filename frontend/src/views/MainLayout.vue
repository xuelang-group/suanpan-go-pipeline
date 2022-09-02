<template>
  <div class="main-layout">
    <div class="main-top">
      <slot name="top"></slot>
    </div>
    <div 
      class="main-layout-left"
      :style="{ width: `${leftPanelWidth}px`, zIndex: moveDirection == 'left' ? 9: 1 }">
       <slot name="left"></slot>
      <div class="main-layout-resize main-layout-resize-left" :class="{ dragging: moveDirection === 'left' }" @mousedown="startMove($event, 'left')"></div>
    </div>
    <div 
      class="main-layout-middle"
      :style="{ marginLeft: `${leftPanelWidth}px`, marginRight: `${showRight ? rightPanelWidth : 0}px` }"
      >
        <slot name="middle"></slot>
      </div>
    <!-- <div
      class="main-layout-bottom"
      :style="{ left: `${leftPanelWidth}px`, right: `${showRight ? rightPanelWidth : 0}px`, height: `${bottomPanelHeight}px` }"
      >
      <slot name="bottom"></slot>
      <div 
        class="main-layout-resize main-layout-resize-bottom" 
        :class="{ dragging: moveDirection === 'bottom' }" @mousedown="startMove($event, 'bottom')"></div>
    </div> -->
    <div 
      class="main-layout-right"
      :style="{ width: `${showRight ? rightPanelWidth : 0}px`, zIndex: moveDirection == 'right' ? 9: 1 }"
      >
      <slot name="right"></slot>
      <div v-show="showRight" class="main-layout-resize main-layout-resize-right" :class="{ dragging: moveDirection === 'right' }" @mousedown="startMove($event, 'right')"></div>
    </div>
  </div>
</template>

<script>

const panelDefaultWidth = 260

export default {
  name: 'MainLayout',
  props: {
    showRight: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      leftPanelWidth: panelDefaultWidth,
      rightPanelWidth: panelDefaultWidth,
      bottomPanelHeight: 40,
      bottomPanelExpand: false,
      moveDirection: null,
    }
  },
  methods: {
    startMove(e, direction) {
      this.moveDirection = direction
      window.addEventListener('mousemove', this.mousemoveHandler)
      window.addEventListener('mouseup', this.mouseupHandler)
    },
    mousemoveHandler(e) {
      e.preventDefault()
      if(!this.moveDirection) {
        return;
      }
      if(this.moveDirection === 'left') {
        this.leftPanelWidth = this.widthRangeRestrict(e.clientX, panelDefaultWidth, window.innerWidth * 0.5)
      }else if(this.moveDirection === 'right') {
        this.rightPanelWidth = this.widthRangeRestrict(window.innerWidth - e.clientX, panelDefaultWidth, window.innerWidth * 0.5)
      }else if(this.moveDirection === 'bottom') {
        this.bottomPanelHeight = this.widthRangeRestrict(window.innerHeight - e.clientY, 40, window.innerHeight * 0.8)
      }
    },
    mouseupHandler(e) {
      this.moveDirection = null
      window.removeEventListener('mousemove', this.mousemoveHandler)
      window.removeEventListener('mouseup', this.mouseupHandler)
    },
    widthRangeRestrict(panelWidth, minVal, maxVal) {
      if(panelWidth < minVal) {
        panelWidth = minVal
      }
      if(panelWidth > maxVal) {
        panelWidth = maxVal
      }
      return panelWidth;
    }
  }
}
</script>

<style>

</style>
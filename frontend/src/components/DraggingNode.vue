<template>
  <div 
    class="dragging-node"
    :style="{ left: `${left}px`, top: `${top}px` }">
    <div class="dragging-node-rect" :style="{ width: `${nodeWidth}px`, height: `${nodeHeight}px` }"></div>
    <div class="dragging-node-icon iconfont">&#xe94b;</div>
    <div class="dragging-node-label">{{ data && data.name }}</div>
    <template v-if="data && data.ports && data.ports.in">
      <div
        v-for="(port, idx) in data.ports.in"
        :key="port.id"
        :style="{ left: `${nodeWidth / (data.ports.in.length + 1) * (idx + 1)}px`, top: '-6px'}"
        class="dragging-node-port">
      </div>
    </template>
    <template v-if="data && data.ports && data.ports.out">
      <div
        v-for="(port, idx) in data.ports.out"
        :key="port.id"
        :style="{ left: `${nodeWidth / (data.ports.out.length + 1) * (idx + 1)}px`, bottom: '-6px'}"
        class="dragging-node-port">
      </div>
    </template>
  </div>
</template>

<script>
import { NodeWidth, NodeHeight } from '../utils/graph'

export default {
  name: 'DraggingNode',
  props: {
    x: {
      type: Number,
      default: 0
    },
    y: {
      type: Number,
      default: 0
    },
    data: {
      type: [Object, null],
      required: true
    }
  },
  data() {
    return {
      nodeWidth: NodeWidth,
      nodeHeight: NodeHeight
    }
  },
  computed: {
    left() {
      return this.x - NodeWidth * 0.5
    },
    top() {
      return this.y - NodeHeight * 0.5
    }
  }
}
</script>

<style>

</style>
<template>
  <g :transform="`translate(${nodeData.x}, ${nodeData.y})`"  @mousedown.stop="nodeSelectHandler">
    <rect 
      :class="['ml-node-border', {'selected-entity':nodeData.ui.selected,'matching-entity':nodeData.ui.matchType === 0,'matched-entity':nodeData.ui.matchType === 1,'port-dragging':portDragging}]" 
      :width="nodeWidth" :height="nodeHeight" rx="8" ry="8" fill="#fff"
     ></rect>
    <text class="ml-graph-icon iconfont" x="10" y="30">&#xe94b;</text>
    <text class="ml-graph-label" x="40" y="26">{{ nodeData.displayName }}<title>{{ nodeData.name }}</title></text>
    <text v-show="!nodeData.ui.valid" class="ml-graph-status-icon iconfont error" text-anchor="end" :x="nodeWidth-10" y="30">&#xe6b8;<title>节点参数不满足要求</title></text>
    <text v-show="nodeData.ui.valid && runningStatus" class="ml-graph-status-icon iconfont success" text-anchor="end" :x="nodeWidth-10" y="30">&#xe6f2;</text>
    <template v-if="nodeData.ports.in.length > 0">
      <g v-for="port in nodeData.ports.in" :key="port.id">
        <circle 
          @mousedown.stop="portDragStart($event, port)" 
          :class="['ml-port-border', 'ml-port-in' , {'port-matching': port.matchingType === 0}, {'port-matched': port.matched }]" 
          :cx="port.x" 
          :cy="port.y"
          r="6">
          <title>{{ port.name }}</title>
        </circle>
        <polygon v-if="port.matched" :class="['ml-port-arrow', {'port-matching': port.matchingType === 0}]" :points="`${port.x-6},${port.y} ${port.x},${port.y+6} ${port.x+6},${port.y}`" />
      </g>
    </template>
    <template v-if="nodeData.ports.out.length > 0">
      <circle 
        v-for="port in nodeData.ports.out" 
        @mousedown.stop="portDragStart($event, port)" 
        :class="['ml-port-border', 'ml-port-out', {'port-matching': port.matchingType === 0}, {'port-matched': port.matched }]" 
        :key="port.id"
        :cx="port.x" 
        :cy="port.y"
        r="6">
        <title>{{ port.name }}</title>
      </circle>
    </template>
  </g>
</template>

<script>
import { NodeWidth, NodeHeight } from '../../utils/graph'

export default {
  name: 'NodeEntity',
  props: {
    runningStatus: {
      type: Boolean
    },
    nodeData: {
      type: Object,
      required: true
    },
    portDragging: {
      type: Boolean,
      default: false
    }
  },
  emits: ['node-select','drag-start', 'port-drag'],
  data() {
    return {
      nodeWidth: NodeWidth,
      nodeHeight: NodeHeight,
    }
  },
  created() {
  },
  beforeUnmount() {
  },
  methods: {
    nodeSelectHandler(e) {
      this.$emit('node-select', this.nodeData, e.ctrlKey, this)
      this.$emit('drag-start', e, this.nodeData, this)
    },
    portDragStart(e, port) {
      if(this.runningStatus) {
        return;
      }
      this.$emit('port-drag', e, port, this.nodeData)
    }
  }
}
</script>

<style>

</style>
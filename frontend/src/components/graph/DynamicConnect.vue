<template>
  <g>
    <!-- <line v-show="!matched" stroke="#777" stroke-width="1"
      :x1="startPosition.x" 
      :y1="startPosition.y"
      :x2="endPosition.x" 
      :y2="endPosition.y" /> -->
    <path :d="d" stroke="#777" stroke-width="1" fill="transparent"></path>
  </g>
</template>

<script>
export default {
  name: 'DynamicConnect',
  props: {
    matched: {
      type: Boolean,
      default: false
    },
    startPosition: {
      type: Object,
    },
    endPosition: {
      type: Object,
    },
    startPort: {
      type: Object
    }
  },
  data() {
    return {}
  },
  computed: {
    d() {
      if(this.matched) {
        if(this.startPort && this.startPort.id.startsWith('out')) {
          return `M ${this.startPosition.x} ${this.startPosition.y} C ${this.startPosition.x} ${this.startPosition.y + 60}, ${this.endPosition.x} ${this.endPosition.y - 60}, ${this.endPosition.x} ${this.endPosition.y}`
        }else {
          return `M ${this.startPosition.x} ${this.startPosition.y} C ${this.startPosition.x} ${this.startPosition.y - 60}, ${this.endPosition.x} ${this.endPosition.y + 60}, ${this.endPosition.x} ${this.endPosition.y}`
        }
      }else {
        return `M ${this.startPosition.x} ${this.startPosition.y} L ${this.endPosition.x} ${this.endPosition.y}`
      }
    }
  }
}
</script>

<style>

</style>
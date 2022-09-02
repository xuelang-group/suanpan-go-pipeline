<template>
  <g>
    <path :d="d" :class="['ml-path', { selected: data.selected }]"></path>
    <path :d="d" class="ml-path-hit" @mousedown.stop="pathSelect"></path>
  </g>
</template>

<script>
export default {
  name: 'ConnectEntity',
  props: {
    data: {
      type: Object
    }
  },
  emits: ['path-select'],
  data() {
    return {}
  },
  computed: {
    d() {
      let sx = this.data.outPort.svgX, sy = this.data.outPort.svgY,
        ex = this.data.inPort.svgX, ey = this.data.inPort.svgY;
      return `M ${sx} ${sy} C ${sx} ${sy + 60}, ${ex} ${ey - 60}, ${ex} ${ey}`;
    }
  },
  methods: {
    pathSelect(e) {
      if(e.buttons === 1) {
        this.$emit('path-select', this.data, e.ctrlKey)
      }
    }
  }
}
</script>

<style>

</style>
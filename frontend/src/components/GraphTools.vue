<template>
  <div>
    <div class="ml-tools ml-tools-slider">
        <span class="iconfont icon-minus" @click="scaleDown"></span>
        <a-slider 
          style="width: 150px"
          v-model:value="sliderScale"
          :marks="sliderScaleMarks"
          :min="sliderMin"
          :max="sliderMax"
          :step="null"
          @change="sliderScaleChanged" />
        <span class="iconfont icon-addition" @click="scaleUp"></span>
    </div>
    <div class="ml-tools ml-tools-clear">
      <a-popconfirm placement="top" ok-text="确定" cancel-text="取消" @confirm="clean">
        <template #title>
          <p>确定清空画布吗？</p>
        </template>
        <a-tooltip placement="top">
        <template #title>
          <span>清空画布</span>
        </template>
        <span class="iconfont icon-clear"></span>
      </a-tooltip>
      </a-popconfirm>
    </div>
  </div>
</template>

<script>
export default {
  name: 'GraphTools',
  props: {
    scale: {
      type: Number,
      required: true
    },
    scaleMarks: {
      type: Array,
    }
  },
  emits: ['scale-to'],
  data() {
    return {
      sliderScale: this.scale
    }
  },
  watch: {
    scale() {
      this.sliderScale = this.scale
    }
  },
  computed: {
    sliderScaleMarks() {
      let scales = {}
      this.scaleMarks.forEach(s => {
        scales[s] = s
      })
      return scales;
    },
    sliderMin() {
      return this.scaleMarks[0]
    },
    sliderMax() {
      return this.scaleMarks[this.scaleMarks.length - 1]
    }
  },
  methods: {
    scaleDown() {
      let idx = this.scaleMarks.findIndex(mark => mark == this.sliderScale)
      idx--
      if(idx < 0) {
        idx = 0
      }
      this.$emit('tool-action', 'scale', this.scaleMarks[idx])
    },
    scaleUp() {
      let idx = this.scaleMarks.findIndex(mark => mark == this.sliderScale)
      idx++
      if(idx > this.scaleMarks.length - 1) {
        idx = this.scaleMarks.length - 1
      }
      this.$emit('tool-action', 'scale', this.scaleMarks[idx])
    },
    sliderScaleChanged() {
      this.$emit('tool-action', 'scale', this.sliderScale)
    },
    clean() {
      this.$emit('tool-action', 'clean')
    }
  }

}
</script>

<style lang="less">
</style>
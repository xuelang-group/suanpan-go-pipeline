<template>
<div class="param-panel-form-vertical">
  <label class="param-panel-form-label">{{ paramInfo.name }}：</label>
   <div class="param-panel-group-inner">
    <a-input-number 
      style="width:100%" 
      v-model:value="currentValue"
      :disabled="readonly"
      @change="changeHandler" />
   </div>
</div>
<div v-show="showError" class="param-validate-tip">{{ errorMsg }}</div>
</template>

<script>
import { requiredValidate, rangeValidate, outRangeMsg } from '../../utils/validate'

export default {
  name: 'MLInputFloat',
  props: {
    paramInfo: {
      type: Object,
      required: true
    },
    params: {
      type: Array,
      required: true
    },
    readonly: {
      type: Boolean
    }
  },
  emits: ['param-change', 'valid-change'],
  data() {
    return {
      currentValue: this.paramInfo.value,
      showError: !this.paramInfo.valid,
      errorMsg: this.paramInfo.errMsg,
    }
  },
  created() {
    this.$parent.paramCompInsts.push(this); 
  },
  beforeUnmount() {
    if(this.$parent.paramCompInsts) {
      let idx = this.$parent.paramCompInsts.indexOf(this)
      if(idx > -1) {
        this.$parent.paramCompInsts.splice(idx, 1)
      }
    }
  },
  watch: {
    showError(val) {
      this.paramInfo.valid = !val
      this.$emit('valid-change')
    },
    errorMsg(val) {
      this.paramInfo.errMsg = val
    }
  },
  mounted() {
  },
  methods: {
    checkValid(val) {
      if(!requiredValidate(val, this.paramInfo, this.params)) {
        this.showError = true
        this.errorMsg = '该参数必填'
      }else if(!rangeValidate(val, this.paramInfo)){
        this.showError = true
        this.errorMsg = outRangeMsg(this.paramInfo)
      }else {
        this.showError = false
        this.errorMsg = ''
      }
    },
    triggerValidCheck() {
      this.checkValid(this.currentValue)
    },
    changeHandler() {
      this.checkValid(this.currentValue)
      this.$emit('param-change', this.currentValue, this.paramInfo, this)
    }
  }
}
</script>

<style>

</style>
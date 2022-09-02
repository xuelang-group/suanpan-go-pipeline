<template>
<div class="param-panel-form-vertical">
  <label class="param-panel-form-label">{{ paramInfo.name }}：</label>
   <div class="param-panel-group-inner">   
    <a-input 
      v-model:value="currentValue" 
      placeholder="数组类型请用逗号分隔"
      :disabled="readonly"
      @change="changeHandler" />
   </div>
</div>
<div v-show="showError" class="param-validate-tip">{{ errorMsg }}</div>
</template>

<script>
import { requiredValidate } from '../../utils/validate'

export default {
  name: 'MLInputArray',
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
      currentValue: this.paramInfo.value ? this.paramInfo.value.join(',') : '',
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
  methods: {
    checkValid(val) {
      if(!requiredValidate(val, this.paramInfo, this.params)) {
        this.showError = true
        this.errorMsg = '该参数必填'
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
      this.$emit('param-change', this.currentValue ? this.currentValue.split(',') : [], this.paramInfo, this)
    }
  }
}
</script>

<style>

</style>
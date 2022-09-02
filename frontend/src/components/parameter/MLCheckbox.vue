<template>
<div class="param-panel-form-horizontal">
  <label class="param-panel-form-label">{{ paramInfo.name }}：</label>
   <div class="param-panel-group-inner">
    <a-checkbox 
      v-model:checked="currentValue"
      :disabled="readonly"
      @change="changeHandler"
    ></a-checkbox>
   </div>
</div>
<div v-show="showError" class="param-validate-tip">{{ errorMsg }}</div>
</template>

<script>
import { requiredValidate } from '../../utils/validate'

export default {
  name: 'MLCheckbox',
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
  watch: {
    showError(val) {
      this.paramInfo.valid = !val
      this.$emit('valid-change')
    },
    errorMsg(val) {
      this.paramInfo.errMsg = val
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
      this.$emit('param-change', this.currentValue, this.paramInfo, this)
    }
  }
}
</script>

<style>

</style>
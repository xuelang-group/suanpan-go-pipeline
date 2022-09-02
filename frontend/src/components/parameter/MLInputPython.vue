<template>
<div class="param-panel-form-vertical">
  <label class="param-panel-form-label">{{ paramInfo.name }}：</label>
   <div class="param-panel-group-inner">
    <div ref="py" class="param-python-input" v-once></div>
   </div>
</div>
<div v-show="showError" class="param-validate-tip">{{ errorMsg }}</div>
</template>

<script>
import CodeMirror from 'codemirror'
import 'codemirror/lib/codemirror.css'
import 'codemirror/theme/eclipse.css'
import 'codemirror/mode/python/python.js'
import 'codemirror/addon/edit/matchbrackets.js'
import 'codemirror/addon/selection/active-line.js'

import { requiredValidate } from '../../utils/validate'

export default {
  name: 'MLInputPython',
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
    },
    readonly() {
      if(this.editor) {
        this.editor.setOption('readOnly', this.readonly ? 'nocursor' : false)
      }
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
  mounted() {
    this.editor = CodeMirror(this.$refs.py, {
      value: this.currentValue || '',
      mode:  "python",
      theme: "eclipse",
      lineNumbers: true,
      matchBrackets: true,
      styleActiveLine: true,
      indentUnit: 4,
      readOnly: this.readonly ? 'nocursor' : false
    })
    this.editor.on('change', (editor, { text }) => {
      this.currentValue = editor.getValue()
      this.changeHandler()
    })
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
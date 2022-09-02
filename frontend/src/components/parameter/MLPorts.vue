<template>
<div class="param-panel-form-vertical">
  <label class="param-panel-form-label">{{ paramInfo.name }}：</label>
   <div class="param-panel-group-inner">
    <div class="ml-param-ports">
      <div 
        v-for="(port, idx) in currentValue"
        :key="port.id"
        class="ml-param-port"
      >
        <span>{{ port.name }}</span>
        <span v-if="!readonly && (idx === currentValue.length - 1)" class="iconfont icon-node-delete deletable" @click="deleteLastPort"></span>
      </div>
      <div v-if="!readonly" class="ml-param-port ml-param-port-add" @click="addPort">
        <span class="iconfont icon-port-add"></span>
      </div>
    </div>
   </div>
</div>
</template>

<script>
import { deepCopy } from '../../utils'
import * as GraphUtils  from '../..//utils/graph'
import { requiredValidate } from '../../utils/validate'

export default {
  name: 'MLPorts',
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
      currentValue: deepCopy(this.paramInfo.value),
      showError: !this.paramInfo.valid,
      errorMsg: this.paramInfo.errMsg,
    }
  },
  computed: {
    portType() {
      if(this.paramInfo.key === 'inPorts') {
        return 'in'
      }else {
        return 'out'
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
    },
    addPort() {
      let portId = this.portType === 'in' ? `in${this.currentValue.length + 1}` : `out${this.currentValue.length + 1}`
      let port = {
        id: portId,
        name: portId
      }
      this.currentValue.push({
        id: portId,
        name: portId
      })
      GraphUtils.addPort(this.portType, this.$store.state.graph.selNodeDatas[0], port)
      this.changeHandler()
    },
    deleteLastPort() {
      this.currentValue.splice(this.currentValue.length - 1, 1)
      GraphUtils.deleteLastPort(this.portType, this.$store.state.graph.selNodeDatas[0], this.$store.state.graph.connectDatas)
      this.changeHandler()
    }
  }
}
</script>

<style>

</style>
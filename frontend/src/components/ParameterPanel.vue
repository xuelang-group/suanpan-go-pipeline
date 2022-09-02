<template>
  <div class="param-panel">
    <template v-if="selNodeDatas.length === 0"></template>
    <template v-if="selNodeDatas.length === 1">
      <div class="param-panel-title">{{ selNodeDatas[0].name }}</div>
      <MLUUID style="margin-bottom: 20px" :uuid="selNodeDatas[0].uuid"></MLUUID>
      <div class="param-panel-form" :key="selNodeDatas[0].uuid">
        <div 
          v-for="param in selNodeDatas[0].parameters"
          :key="param.key"
          class="param-panel-form-group">
          <component
            :is="paramTypes[param.type]"
            :params="selNodeDatas[0].parameters"
            :param-info="param"
            :readonly="!this.$store.getters['graph/editable']"
            @param-change="paramChanged"
            @valid-change="validChanged">
          </component>
        </div>
      </div>
    </template>
    <template v-if="selNodeDatas.length > 1">
      <div 
        v-for="selNodeData in selNodeDatas"
        :key="selNodeData.uuid"
        class="param-panel-multi"
        >
        <div class="param-panel-multi-row">
          <div class="param-panel-multi-label">UUID</div>
          <div class="param-panel-multi-content">{{ selNodeData.uuid }}</div>
        </div>
        <div class="param-panel-multi-row">
          <div class="param-panel-multi-label">节点名称</div>
          <div class="param-panel-multi-content">{{ selNodeData.name }}</div>
        </div>
      </div>
    </template>
  </div>
</template>

<script>
import MLInputString from './parameter/MLInputString.vue'
import MLSelect from './parameter/MLSelect.vue'
import MLInputFloat from './parameter/MLInputFloat.vue'
import MLInteger from './parameter/MLInteger.vue'
import MLInputArray from './parameter/MLInputArray.vue'
import MLGroupSelect from './parameter/MLGroupSelect.vue'
import MLUUID from './parameter/MLUUID.vue'
import MLFileUploader from './parameter/MLFileUploader.vue'
import MLInputPython from './parameter/MLInputPython.vue'
import MLPorts from './parameter/MLPorts.vue'
import MLCheckbox from './parameter/MLCheckbox.vue'
import { getNodeValidation } from '../utils/graph'

import { markRaw } from 'vue'

export default {
  name: 'ParameterPanel',
  components: {
    MLInputString,
    MLSelect,
    MLInputFloat,
    MLInteger,
    MLInputArray,
    MLGroupSelect,
    MLUUID,
    MLFileUploader,
    MLInputPython,
    MLPorts,
    MLCheckbox
  },
  data() {
    return {
      paramTypes: {
        'inputString': 'MLInputString',
        'inputFloat': 'MLInputFloat',
        'select': 'MLSelect',
        'inputArray': 'MLInputArray',
        'inputInteger': 'MLInteger',
        'groupSelect': 'MLGroupSelect',
        'fileUploader': 'MLFileUploader',
        'inputPythonScript': 'MLInputPython',
        'ports': 'MLPorts',
        'checkbox': 'MLCheckbox'
      },
      paramCompInsts: markRaw([])
    }
  },
  computed: {
    selNodeDatas() {
      return this.$store.state.graph.selNodeDatas
    }
  },
  beforeUnmount() {
    this.paramCompInsts = null
  },
  methods: {
    paramChanged(val, param, compInst) {
      param.value = val
      this.paramCompInsts.forEach(inst => {
        if(inst !== compInst) {
          inst.triggerValidCheck()
        }
      })
      this.$store.dispatch('graph/update')
    },
    validChanged() {
      this.selNodeDatas[0].ui.valid = getNodeValidation(this.selNodeDatas[0].parameters)
    }
  }
}
</script>

<style>

</style>
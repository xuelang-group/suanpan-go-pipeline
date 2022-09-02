<template>
  <div class="ml-top-header">
    <div class="ml-graph-action">
      <MLButton 
        :iconfont="$store.state.graph.graphStatus === 1 ? 'icon-editor' : 'icon-deploy'" 
        :label="$store.state.graph.graphStatus === 1 ? '编辑': '部署'" 
        :disabled="$store.state.graph.processStatus === 1"
        @click="changeGraphStatus"
      >
      </MLButton>
      <!-- <div 
        class="ml-graph-action-item"
        :class="{ disabled: $store.state.graph.processStatus === 1 }"
        @click="changeGraphStatus">
         <span :class="['iconfont', {'icon-editor': $store.state.graph.graphStatus === 1}, {'icon-deploy': $store.state.graph.graphStatus === 0}]"></span>
         <label>{{ $store.state.graph.graphStatus === 1 ? '编辑': '部署' }}</label>
      </div> -->
      <!-- <div class="ml-graph-action-item" @click="saveGraphModel">
         <span class="iconfont icon-save"></span>
         <label>发布</label>
      </div>
      <div 
        class="ml-graph-action-item" 
        :class="{ disabled: $store.state.graph.graphStatus === 1 || $store.state.graph.processStatus === 1 }" 
        @click="processHandler">
         <span class="iconfont icon-run"></span>
         <label>运行</label>
      </div> -->
    </div>
    <GraphModelPublish v-if="modelCreateModal" v-model:visible="modelCreateModal" @confirm="createModel"></GraphModelPublish>
  </div>
</template>

<script>
import { message } from 'ant-design-vue'
import { updateGraphStatus, publishModel, runProcess } from '../service'
import * as GraphUtils  from '../utils/graph'
import GraphModelPublish from './GraphModelPublish.vue'
import MLButton from './common/MLButton.vue'
import { EVENT_MODEL_PUBLISH } from '../utils/event'
import BUS from '../utils/bus'

export default {
  name: 'TopHeader',
  components: {
    GraphModelPublish,
    MLButton
  },
  data() {
    return {
      modelCreateModal: false
    }
  },
  methods: {
    changeGraphStatus() {
      if(this.$store.state.graph.processStatus === 1) {
        return
      }
      if(this.$store.state.graph.graphStatus === 0) {
        // check parameter valid
        this.$store.dispatch('graph/validateGraph').then(() => {
          this.$store.commit('graph/graphStatus', 1)
          updateGraphStatus({status: 1, graph: GraphUtils.toGraphRawData(this.$store.state.graph)})
            .catch(err => {
              console.error('update graph status error:', err)
            })
        }).catch(errMsg => {
          message.warning(errMsg);
        })
      }else {
        this.$store.commit('graph/graphStatus', 0)
        updateGraphStatus({status: 0})
          .catch(err => {
            console.error('update graph status error:', err)
          })
      }
    },
    saveGraphModel() {
      this.modelCreateModal = true
    },
    createModel(modelState) {
      modelState.graph = GraphUtils.toGraphRawData(this.$store.state.graph)
      publishModel(modelState).then((model) => {
        BUS.emit(EVENT_MODEL_PUBLISH, [model])
        message.success(`模型"${model.name}"发布成功`)
      }).catch(err => {
        console.error('save model error:', err)
        message.error(`模型"${model.name}"发布失败`)
      })
    },
    processHandler() {
      if(this.$store.state.graph.graphStatus === 1) {
        // 部署状态无效
        return;
      }
      if(this.$store.state.graph.processStatus === 1) {
        return;
      }
      if(this.$store.state.graph.processStatus === 0) {
        // check parameter valid
        this.$store.dispatch('graph/validateGraph').then(() => {
          this.$store.commit('graph/processStatus', 1)
          runProcess(GraphUtils.toGraphRawData(this.$store.state.graph))
        }).catch(err => {
          message.error(err)
          console.error(err)
        })
      }
    }
  }
}
</script>

<style>

</style>
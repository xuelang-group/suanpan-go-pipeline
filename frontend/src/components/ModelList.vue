<template>
  <div class="ml-model-container">
    <div ref="container" class="ml-model-list-view" @scroll="handleScroll">
      <div 
        class="ml-model-phantom"
        :style="{ height: `${itemHeight * modelList.length}px` }">
      </div>
       <div
        ref="content"
        class="ml-model-list-content">
        <div
          v-for="modelItem in visibleModelList"
          :key="modelItem.id"
          class="ml-model-list-item"
          :style="{
            height: `${itemHeight}px`
          }"
          >
          <span class="ml-model-item-label" :title="modelItem.name">{{ modelItem.name }}</span>
          <div class="ml-model-item-btns">
            <span class="ml-model-item-btn update" @click="updateModel(modelItem)">修改</span>
            <a-popconfirm
              title="确定删除该模型吗？"
              ok-text="确定"
              cancel-text="取消"
              @confirm="deleteModel(modelItem)"
            >
              <span class="ml-model-item-btn delete">删除</span>
            </a-popconfirm>
            <span 
              class="ml-model-item-btn check" 
              :class="{'disabled': !$store.getters['graph/editable']}"
              @click="modelView(modelItem)">查看</span>
          </div>
      </div>
    </div>
    </div>
    <GraphModelUpdate v-if="modelUpdateModal" v-model:visible="modelUpdateModal" :model="updatingModel" @confirm="updateModelHandler"></GraphModelUpdate>
    <div v-if="loading" class="ml-loading">
      <a-spin />
    </div>
  </div>
</template>

<script>
import { getModelList, 
  deleteModel as deleteModelService, 
  updateModel as updateModelService,
  selectModel as selectModelService } from '../service/'
import BUS from '../utils/bus'
import { EVENT_MODEL_PUBLISH } from '../utils/event'
import { message, Modal } from 'ant-design-vue'
import GraphModelUpdate from './GraphModelUpdate.vue'

export default {
  name: 'ModelList',
  components: {
    GraphModelUpdate
  },
  data() {
    return {
      itemHeight: 32,
      modelList: [],
      visibleModelList: [],
      loading: false,
      modelUpdateModal: false,
      updatingModel: null
    }
  },
  created() {
    BUS.on(EVENT_MODEL_PUBLISH, this.modelSave)
  },
  mounted() {
    this.getModels();
    this.resizeObserver = new ResizeObserver( () => {
      this.handleScroll()
    })
    this.resizeObserver.observe(this.$refs.container)
  },
  beforeUnmount() {
    this.resizeObserver.disconnect()
    BUS.off(EVENT_MODEL_PUBLISH, this.modelSave)
  },
  methods: {
    getModels() {
      this.loading = true
      getModelList().then(data => {
        this.modelList = data
        this.handleScroll()
      }).catch(err => {
        console.error('get model list error:', err)
      }).finally(() => {
        this.loading = false
      })
    },

    updateVisibleData(scrollTop) {
      scrollTop = scrollTop || 0;
      const visibleCount = Math.ceil(this.$refs.container.clientHeight / this.itemHeight);
      const start = Math.floor(scrollTop / this.itemHeight);
      const end = start + visibleCount + 1;
      this.visibleModelList = this.modelList.slice(start, end);
      this.$refs.content.style.webkitTransform = `translate3d(0, ${ start * this.itemHeight }px, 0)`;
    },

    handleScroll() {
      const scrollTop = this.$refs.container.scrollTop;
      this.updateVisibleData(scrollTop);
    },

    modelSave(newModel) {
      let idx = this.modelList.findIndex(model => model.id === newModel.id)
      if(idx > -1) {
        this.modelList.splice(idx, 1, newModel)
      }else {
        this.modelList.push(newModel)
      }
      this.handleScroll()
    },

    deleteModel(model) {
      deleteModelService({id: model.id}).then(() => {
        message.success(`模型"${model.name}"删除成功`)
        let idx = this.modelList.findIndex(m => m.id === model.id)
        if(idx > -1) {
          this.modelList.splice(idx, 1)
        }
        this.handleScroll()
      }).catch(err => {
        console.log('delete model error:', err)
        message.error(`模型"${model.name}"删除失败`)
      })
    },

    updateModel(modelItem) {
      this.modelUpdateModal = true
      this.updatingModel = modelItem
    },

    updateModelHandler(newModel) {
      updateModelService(newModel).then(data => {
        message.success(`更新成功`)
        let idx = this.modelList.findIndex(m => m.id === data.id)
        if(idx > -1) {
          this.modelList.splice(idx, 1, data)
        }
        this.handleScroll()
      }).catch(err => {
        message.error(`更新失败`)
        console.error('model upate error:', err)
      })
    },

    modelView(modelItem) {
      if(!this.$store.getters['graph/editable']) {
        return
      }
      if(!this.checkGraph()) {
         Modal.warning({
            title: () => '警告',
            content: () => `模型"${modelItem.name}"将覆盖当前图的所有信息，是否继续？`,
            okText: '继续',
            closable: true,
            onOk: () => {
              this.selectModel(modelItem)
            }
          })
      }else {
        this.selectModel(modelItem)
      }
    },

    selectModel(modelItem) {
      this.$store.commit('graph/graphLoading', true)
      selectModelService({id: modelItem.id}).then(graphData => {
          this.$store.dispatch('graph/generateGraph', graphData).then( () => {
            this.$store.dispatch('graph/update')
          })
        }).catch(err => {
          console.error('model select error:', err)
        }).finally(() => {
          this.$store.commit('graph/graphLoading', false)
        })
    },

    checkGraph() {
      if(this.$store.state.graph.nodeDatas.length < 1) {
        return true
      }else {
        return false
      }
    }
  }
}
</script>

<style>

</style>
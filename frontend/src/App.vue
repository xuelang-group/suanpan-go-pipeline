<template>
<div>
  <MainLayout :show-right="$store.state.graph.selNodeDatas.length > 0">
    <template v-slot:top>
      <TopHeader></TopHeader>
    </template>
    <template v-slot:middle>
      <GraphEditor ref="graph"></GraphEditor>
    </template>
    <template v-slot:left>
      <MLCollapse>
        <MLCollapsePanel 
          header="组件列表"
          :expend="true">
          <ComponentList @start-drag="startDrag" @dragging="dragging" @end-drag="endDrag"></ComponentList>
        </MLCollapsePanel>
        <!-- <MLCollapsePanel 
          header="流程列表"
          :expend="true">
          <ModelList></ModelList>
        </MLCollapsePanel> -->
      </MLCollapse>
    </template>
    <template v-slot:right>
      <ParameterPanel />
    </template>
  </MainLayout>
  <DraggingNode 
    v-show="draggingVisible"
    :x="draggingX"
    :y="draggingY"
    :data="draggingNode"
    :style="{ cursor: dropAllowed ? 'default' : 'not-allowed' }"
     />
  <InitialLoading v-if="$store.state.initLoading" />
</div>
</template>

<script>
import MainLayout from './views/MainLayout.vue'
import MLCollapse from './components/collapse/MLCollapse.vue'
import MLCollapsePanel from './components/collapse/MLCollapsePanel.vue'
import ComponentList from './components/ComponentList.vue'
import TopHeader from './components/TopHeader.vue'
import GraphEditor from './components/GraphEditor.vue'
import DraggingNode from './components/DraggingNode.vue'
import ParameterPanel from './components/ParameterPanel.vue'
import InitialLoading from './components/InitialLoading.vue'
import ModelList from './components/ModelList.vue'
import SocketIOClient from './service/sio'

import { message } from 'ant-design-vue'
import { deepCopy }  from './utils/index'
import { getComponentList, getGraph, getGraphStatus, 
  getProcessStatus, ossServiceInit, getStorageInfo } from './service'
import * as GraphUtils  from './utils/graph'

export default {
  name: 'App',
  components: {
    MainLayout,
    MLCollapse,
    MLCollapsePanel,
    ComponentList,
    TopHeader,
    GraphEditor,
    DraggingNode,
    ParameterPanel,
    InitialLoading,
    ModelList
  },
  data() {
    return {
      draggingVisible: false,
      draggingX: 0,
      draggingY: 0,
      draggingNode: null,
      dropAllowed: false
    }
  },
  created() {
    SocketIOClient.on('process.status', this.processStatusHandler)
  },
  mounted() {
    Promise.all([
        getComponentList(), 
        getGraph(), 
        getGraphStatus(), 
        getProcessStatus()
      ]
      ).then(data => {
        this.$store.commit('graph/componentRawData', data[0])
        this.$store.commit('graph/graphStatus', data[2])
        this.$store.commit('graph/processStatus', data[3].status)
        if(data[4] && data[4]['type']) {
          ossServiceInit(data[4]['type'])
          this.$store.commit('storageNodePath', data[4].nodePath)
        }
        this.$store.dispatch('graph/generateGraph', data[1])
        this.$store.commit('initLoading', false)
    }).catch(err => {
      console.error('getComponentList error: ', err)
    })
  },
  methods: {
    startDrag(componentData, x, y) {
      if(!this.$store.getters["graph/editable"]) {
        return
      }
      this.draggingVisible = true
      this.draggingX = x
      this.draggingY = y
      this.draggingNode = componentData
    },
    dragging(x, y) {
      if(!this.draggingVisible) {
        return
      }
      this.draggingX = x
      this.draggingY = y
      if(this.$store.state.graph.graphBounding) {
        if((x > this.$store.state.graph.graphBounding.left) && (x < this.$store.state.graph.graphBounding.right)
          && (y > this.$store.state.graph.graphBounding.top) && (y < this.$store.state.graph.graphBounding.bottom)) {
          this.dropAllowed = true
        }else {
          this.dropAllowed = false
        }
      }
    },
    endDrag(e) {
      if(!this.draggingVisible) {
        return
      }
      this.draggingVisible = false
      if(this.dropAllowed && this.$refs.graph) {
        let res = GraphUtils.validateGraphBeforeDrop(this.draggingNode, this.$store.state.graph.nodeDatas)
        if(res.error) {
          message.warning(res.error.message)
        }else {
          this.$refs.graph.addNode(
              deepCopy(this.draggingNode), 
              e.clientX, 
              e.clientY
            );
          this.draggingNode = null
        }
      }
    },
    processStatusHandler(data) {
      console.log('processStatus:', data)
      this.$store.commit('graph/processStatus', data.status)
    }
  }
}
</script>

<style lang="less">
</style>

<template>
  <div class="ml-editor">
    <div class="ml-editor-tip" 
      v-show="($store.getters['graph/editable']) && ($store.state.graph.graphMode === 2)">
      <span class="iconfont icon-graph-sync"></span>
      <label>同步中...</label>
    </div>
    <div
      ref="graph"
      class="ml-editor-graph"
      @mousedown="cancelSelect"
      @contextmenu="contextmenuHandler">
      <svg 
        :width="graphW" 
        :height="graphH"
        id="ml-graph"
        class="ml-graph"
        xmlns="http://www.w3.org/2000/svg">
        <g :transform="`matrix(${graphScale}, 0, 0, ${graphScale}, ${graphTransX}, ${graphTransY})`">
          
           <StaticConnect
            v-for="connectData in connectDatas"
            :key="connectData.uuid"
            :data="connectData"
            @path-select="connectSelect"
           ></StaticConnect>
           
           <DynamicConnect
            v-show="draggingPortFlg"
            :start-position="draggingPortStart"
            :end-position="draggingPortEnd"
            :matched="!!matchedPort"
            :start-port="draggingPort"
            ></DynamicConnect>
          
          <NodeEntity 
            v-for="nodeData in nodeDatas" 
            :key="nodeData.uuid"
            :running-status="!$store.getters['graph/editable']"
            :node-data="nodeData"
            :port-dragging="draggingPortFlg"
            @node-select="nodeSelect"
            @drag-start="nodeDragStart"
            @port-drag="portDragStart"
            >
          </NodeEntity>

        </g>
      </svg>
      <teleport to="body">
        <GraphContextMenu 
          v-if="contextMenuVisible"
          :context-menus="contextMenus"
          v-model:visible="contextMenuVisible" 
          :mouse-position="contextMenuPos"
          @menu-click="contextMenuClickHandler">
        </GraphContextMenu>
      </teleport>
    </div>
    <GraphTools 
      class="ml-editor-tools"
      :scale="graphScale"
      :scale-marks="graphScales"
      @tool-action="toolActionHandler">
    </GraphTools>
    <GraphNodeData v-if="nodeDataVisible" v-model:visible="nodeDataVisible" :data="$store.state.graph.selNodeDatas"></GraphNodeData>
    <div v-show="$store.state.graph.graphLoading" class="ml-loading">
      <a-spin size="large" />
    </div>
  </div>
</template>

<script>
import NodeEntity from './graph/NodeEntity.vue'
import GraphContextMenu from './GraphContextMenu.vue'
import GraphNodeData from './GraphNodeData.vue'
import DynamicConnect from './graph/DynamicConnect.vue'
import StaticConnect from './graph/StaticConnect.vue'
import GraphTools from './GraphTools.vue'
import * as GraphUtils  from '../utils/graph'
import { debounce }  from '../utils/index'
import { message } from 'ant-design-vue'

export default {
  name: 'GraphEditor',
  components: {
    NodeEntity,
    GraphContextMenu,
    DynamicConnect,
    StaticConnect,
    GraphTools,
    GraphNodeData
  },
  data() {
    return {
      graphW: 0,
      graphH: 0,
      draggingGraph: false,
      contextMenuVisible: false,
      contextMenuPos: null,
      draggingPortFlg: false,
      draggingPortStart: {x:0,y:0},
      draggingPortEnd: {x:0,y:0},
      draggingPort: null,
      matchedPort: null,
      matchedNode: null,
      nodeDataVisible: false,
    }
  },
  computed: {
    contextMenus() {
      return [{
          key: 'delete',
          disabled: (!this.$store.getters["graph/editable"]) 
            || (this.$store.state.graph.selNodeDatas.length === 0 && this.$store.state.graph.selConnectionDatas.length === 0),
          icon: 'icon-node-delete',
          label: '删除'
        }, {
          key: 'copy',
          disabled: (!this.$store.getters["graph/editable"]) 
            || (this.$store.state.graph.selNodeDatas.length === 0 && this.$store.state.graph.selConnectionDatas.length === 0),
          icon: 'icon-node-copy',
          label: '复制'
        }, {
          key: 'paste',
          disabled: !this.$store.getters["graph/editable"],
          icon: 'icon-node-paste',
          label: '粘贴'
        }, {
          key: 'divider1',
          type: 'divider',
          style: { 'pointer-events': 'none' }
        }, {
          key: 'data',
          disabled: !(this.$store.state.graph.selNodeDatas.length === 1 && this.$store.state.graph.selConnectionDatas.length === 0),
          icon: 'icon-node-data-view',
          label: '结果数据',
          children: [{
            key: 'data-view',
            icon: 'icon-node-data-view',
            disabled: false,
            label: '数据预览',
          }, {
            key: 'data-model-download',
            icon: 'icon-node-data-view',
            disabled: false,
            label: '节点模型下载',
          }, {
            key: 'data-model-download',
            icon: 'icon-node-data-view',
            disabled: false,
            label: '下载',
          }]
        }, {
          key: 'divider2',
          type: 'divider',
          style: { 'pointer-events': 'none' }
        }, {
          key: 'help',
          disabled: !(this.$store.state.graph.selNodeDatas.length === 1 && this.$store.state.graph.selConnectionDatas.length === 0),
          icon: 'icon-node-help',
          label: '帮助文档'
        }]
    },
    nodeDatas() {
      return this.$store.state.graph.nodeDatas
    },
    connectDatas() {
      return this.$store.state.graph.connectDatas
    },
    graphScale: {
      get() {
        return this.$store.state.graph.graphScale
      },
      set(val) {
        this.$store.state.graph.graphScale = val
      }
    },
    graphScales() {
      return this.$store.state.graph.graphScales
    },
    graphScaleIndex: {
      get() {
        return this.$store.state.graph.graphScaleIndex
      },
      set(val) {
        this.$store.state.graph.graphScaleIndex = val
      }
    },
    graphTransX: {
      get() {
        return this.$store.state.graph.graphTransX
      },
      set(val) {
        this.$store.state.graph.graphTransX = val
      }
    },
    graphTransY: {
      get() {
        return this.$store.state.graph.graphTransY
      },
      set(val) {
        this.$store.state.graph.graphTransY = val
      }
    }
  },
  created() {
    this.scaleWithPointerDebounce = debounce(this.scaleWithPointer, 80)
  },
  mounted() {
    this.resizeObserver = new ResizeObserver( () => {
      this.updateSvgSize()
    })
    this.resizeObserver.observe(this.$refs.graph)
    window.addEventListener("wheel", this.scaleWithPointerDebounce)
    window.addEventListener("keydown", this.keydownHandler)
  },
  beforeUnmount() {
    this.resizeObserver.disconnect()
    window.removeEventListener("wheel", this.scaleWithPointerDebounce)
    window.removeEventListener("keydown", this.keydownHandler)
  },
  watch: {
    graphScale() {
      this.$store.dispatch('graph/update')
    }
  },
  methods: {
    updateSvgSize() {
      this.graphW = this.$refs.graph.clientWidth
      this.graphH = this.$refs.graph.clientHeight
      this.$store.commit('graph/graphBounding', this.$refs.graph.getBoundingClientRect())
    },
    cancelSelect(e) {
      window.addEventListener('mousemove', this.graphMousemove)
      window.addEventListener('mouseup', this.graphMouseup)

      if(e.buttons === 1) {
        this.nodeDatas.forEach(nodeData => {
          nodeData.ui.selected = false
        })
        this.connectDatas.forEach(connectData => {
          connectData.selected = false
        })
        this.$store.commit('graph/selNodeDatas', [])
        this.$store.commit('graph/selConnectionDatas', [])
      }

      this.draggingGraph = true
      this.dragPrevX = e.clientX
      this.dragPrevY = e.clientY
      this.updateFlag = false
    },
    nodeSelect(nodeData, multiSelect) {
      if(!multiSelect && this.$store.state.graph.selNodeDatas.some(selNodeData => selNodeData.uuid === nodeData.uuid)) {
        multiSelect = true
      }
      this.entitySelect(nodeData, multiSelect)
    },
    graphMousemove(e) {
      e.preventDefault()
      if(!this.draggingGraph) {
        return;
      }
      this.updateFlag = true
      let diffX = e.clientX - this.dragPrevX,
        diffY = e.clientY - this.dragPrevY;
      this.dragPrevX = e.clientX
      this.dragPrevY = e.clientY
      this.graphTransX = this.graphTransX + diffX
      this.graphTransY = this.graphTransY + diffY
    },
    graphMouseup(e) {
      this.draggingGraph = false
      window.removeEventListener('mousemove', this.graphMousemove)
      window.removeEventListener('mouseup', this.graphMouseup)
      if(this.updateFlag) {
        this.$store.dispatch('graph/update')
        this.updateFlag = false
      }
    },
    nodeDragStart(e) {
      this.nodeDraggingFlag = true
      this.nodeDragPrevPos = this.getSvgPos(e.clientX, e.clientY)
      this.draggingNodesPos = this.$store.state.graph.selNodeDatas.map(selNodeData => { return {x:selNodeData.x, y:selNodeData.y }})
      this.updateFlag = false
      window.addEventListener('mousemove', this.nodeDragging)
      window.addEventListener('mouseup', this.nodeDragEnd)
    },
    nodeDragging(e) {
      e.preventDefault()
      if(!this.nodeDraggingFlag) {
        return;
      }
      let pos = this.getSvgPos(e.clientX, e.clientY)
      let diffX = pos.x - this.nodeDragPrevPos.x,
        diffY = pos.y - this.nodeDragPrevPos.y;
      this.$store.state.graph.selNodeDatas.forEach((nodeData, idx) => {
        nodeData.x = this.draggingNodesPos[idx].x + diffX
        nodeData.y = this.draggingNodesPos[idx].y + diffY
        GraphUtils.updatePortsPosition(nodeData)
      })
      this.updateFlag = true
    },
    nodeDragEnd(e) {
      this.nodeDraggingFlag = false
      window.removeEventListener('mousemove', this.nodeDragging)
      window.removeEventListener('mouseup', this.nodeDragEnd)
      if(this.updateFlag) {
        this.$store.dispatch('graph/update')
        this.updateFlag = false
      }
    },
    portDragStart(e, port, nodeData) {
      this.draggingPortFlg = true
      this.draggingPortStart = {x: port.x + nodeData.x, y: port.y + nodeData.y}
      this.draggingPortEnd = this.draggingPortStart
      this.draggingPort = port
      this.draggingPortNode = nodeData
      window.addEventListener('mousemove', this.portDragging)
      window.addEventListener('mouseup', this.portDragEnd)
    },
    portDragging(e) {
      if(!this.draggingPortFlg) {
        return
      }
      let x = e.clientX- this.$store.state.graph.graphBounding.left, 
          y = e.clientY - this.$store.state.graph.graphBounding.top;
      let {x: svgX, y: svgY} = GraphUtils.containerCoordinate2Svg(x, y, this.graphTransX, this.graphTransY, this.graphScale)

      let portsKey
      if(this.draggingPort.id.startsWith('out')) {
        portsKey = 'in'
      }else {
        portsKey = 'out'
      }
      this.nodeDatas.forEach(nodeData => {
        if((nodeData.uuid != this.draggingPortNode.uuid) && (nodeData.ports[portsKey].length > 0)) {
          nodeData.ui.matchType = 0
          nodeData.ports[portsKey].forEach(port => {
            port.matchingType = 0
          })
        }
      })

      // TODO: 可以使用四叉树优化
      // let curNodeDatas = this.nodeDatas.filter(nodeData => {
      //   return GraphUtils.contains(svgX, svgY, nodeData.x, nodeData.y, GraphUtils.NodeWidth, GraphUtils.NodeHeight)
      // })
      this.matchedPort = null
      this.matchedNode = null
      for(let i = this.nodeDatas.length - 1; i > -1; i--) {
        let node = this.nodeDatas[i], ports = null;
        if(node.uuid === this.draggingPortNode.uuid) {
          continue
        }
        if(this.draggingPort.id.startsWith('out')) {
          ports = node.ports.in
        }else{
          ports = node.ports.out
        }
        for(let j = 0, portLen = ports.length; j < portLen; j++) {
          if(GraphUtils.contains2(svgX, svgY, ports[j].svgBLeft, ports[j].svgBRight, ports[j].svgBTop, ports[j].svgBBottom)) {
            this.matchedPort = ports[j]
            this.matchedNode = node
            this.matchedNode.ui.matchType = 1
            break
          }
        }
        if(this.matchedPort) {
          break
        }
      }
      if(this.matchedPort) {
        this.draggingPortEnd = {x: this.matchedPort.svgX, y: this.matchedPort.svgY}
      }else {
        this.draggingPortEnd = {x: svgX, y: svgY}
      }
    
    },
    portDragEnd(e) {
      window.removeEventListener('mousemove', this.portDragging)
      window.removeEventListener('mouseup', this.portDragEnd)

      if(this.matchedPort) {
        let inPort, inNode, outPort, outNode;
        if(this.draggingPort.id.startsWith('out')) {
          inPort = this.matchedPort
          inNode = this.matchedNode
          outPort = this.draggingPort
          outNode = this.draggingPortNode
        }else {
          inPort = this.draggingPort
          inNode = this.draggingPortNode
          outPort = this.matchedPort
          outNode = this.matchedNode
        }
        if(this.checkBeforeAddConnect(inPort, inNode, outPort, outNode)) {
          this.addConnect(inPort, inNode, outPort, outNode);
        }
        this.matchedPort = null
      }

      this.draggingPortFlg = false
      this.draggingPortStart = {x:0,y:0}
      this.draggingPortEnd = {x:0,y:0}
      this.draggingPort = null
      this.draggingPortNode = null
      this.nodeDatas.forEach(nodeData => {
        nodeData.ui.matchType = null
        nodeData.ports.in.forEach(port => {
          port.matchingType = null
        })
        nodeData.ports.out.forEach(port => {
          port.matchingType = null
        })
      })
    },
    addConnect(inPort,inNode,outPort,outNode) {
      inPort.matched = true
      outPort.matched = true
      let connector = GraphUtils.generateConnect(inPort,inNode,outPort,outNode)

      this.connectDatas.push(connector)
      this.$store.dispatch('graph/update')
    },
    deleteConnect(connectData) {
      GraphUtils.deleteConnect(connectData, this.connectDatas)

      this.$store.dispatch('graph/update')
    },
    connectSelect(connectData, multiSelect) {
      this.entitySelect(connectData, multiSelect)
    },
    entitySelect(entityData, multiSelect) {
      this.contextMenuVisible = false
      if(multiSelect) {
        // 多选
        if(entityData.key) {
          // 选中了节点
          entityData.ui.selected = true
          if(!this.$store.state.graph.selNodeDatas.some(selNodeData => selNodeData.uuid === entityData.uuid )) {
            this.$store.state.graph.selNodeDatas.push(entityData)
          }
        }else {
          // 选中了连线
          entityData.selected = true
          if(!this.$store.state.graph.selNodeDatas.some(selNodeData => selNodeData.uuid === entityData.uuid )) {
            this.$store.state.graph.selConnectionDatas.push(entityData)
          }
        }
      }else {
        // 单选
        this.$store.state.graph.selNodeDatas.forEach(nodeData => {
          nodeData.ui.selected = false
        })
        this.$store.state.graph.selConnectionDatas.forEach(connectData => {
          connectData.selected = false
        })
        if(entityData.key) {
          // 选中了节点
          entityData.ui.selected = true
          this.$store.commit('graph/selNodeDatas', [entityData])
          this.$store.commit('graph/selConnectionDatas', [])
        }else {
          // 选中了连线
          entityData.selected = true
          this.$store.commit('graph/selNodeDatas', [])
          this.$store.commit('graph/selConnectionDatas', [entityData])
        }
      }
    },
    checkBeforeAddConnect(inPort, inNode, outPort, outNode) {
      // 检查是否已经存在
      let inConnectors = GraphUtils.nodeGetPortConnectorsByPort(inNode.ui.inLines, 'in', inNode, inPort)
      let outConnectors = GraphUtils.nodeGetPortConnectorsByPort(outNode.ui.outLines, 'out', outNode, outPort)
      
      let inConnectorsSet = new Set(inConnectors.map( c => c.uuid))
      let outConnectorsSet = new Set(outConnectors.map( c => c.uuid))
      let intersection = new Set([...inConnectorsSet].filter(inC => outConnectorsSet.has(inC)));
      
      if(intersection.keys.length > 0) {
        return false;
      }

      // 输入口已经连接过
      if(inConnectors.length > 0) {
        inConnectors.forEach(connector => {
          this.deleteConnect(connector)
        })
      }
      
      return true
    },
    scaleWithPointer(e) {
      const rect = this.$store.state.graph.graphBounding
      let clientX = e.clientX,
        clientY = e.clientY;
      if(clientX < rect.left || clientX > rect.right 
        || clientY < rect.top || clientY > rect.bottom) {
        return;
      }
      let currentPosX = e.clientX - rect.left,
        currentPosY = e.clientY - rect.top;
      let currentScale = this.graphScale;
      let targetScale;
      if(e.deltaY < 0) {
        targetScale = this.scaleUp()
      }else {
        targetScale = this.scaleDown()
      }
      this.scaleToWithPointer(currentPosX, currentPosY, currentScale, targetScale)
    },
    sliderScaleTo(targetScale) {
      let currentScale = this.graphScale
      this.graphScale = targetScale
      const rect = this.$store.state.graph.graphBounding
      this.scaleToWithPointer(rect.width * 0.5, rect.height * 0.5, currentScale, targetScale)
    },
    toolActionHandler(toolAction, param) {
      if(toolAction === 'scale') {
        this.sliderScaleTo(param)
      }else if(toolAction === 'clean') {
        this.$store.commit('graph/clean')
        this.$store.dispatch('graph/update')
      }
    },
    scaleToWithPointer(x, y, currentScale, targetScale) {
      let res = GraphUtils.scaleBypoint(
        x, y,
        this.graphTransX, this.graphTransY, 
        currentScale, targetScale)
      this.graphTransX = res.translateX
      this.graphTransY = res.translateY
    },
    scaleUp() {
      this.graphScaleIndex++
      if(this.graphScaleIndex >= this.graphScales.length) {
        this.graphScaleIndex = this.graphScales.length - 1
      }
      this.graphScale = this.graphScales[this.graphScaleIndex]
      return this.graphScale
    },
    scaleDown() {
      this.graphScaleIndex--
      if(this.graphScaleIndex < 0) {
        this.graphScaleIndex = 0
      }
      this.graphScale = this.graphScales[this.graphScaleIndex]
      return this.graphScale
    },
    contextmenuHandler(e) {
      e.preventDefault()
      this.contextMenuPos = {x: e.clientX, y: e.clientY}
      this.contextMenuVisible = true
    },
    addNode(nodeData, clientX, clientY) {
      this.nodeDatas.push(GraphUtils.generateNode(nodeData, this.getSvgPos(clientX, clientY)))
      this.$store.dispatch('graph/update')
    },
    deleteNode(selNodeData) {
      let idx = this.nodeDatas.findIndex(nodeData => nodeData.uuid === selNodeData.uuid)
      if(idx < 0) {
        return
      }
      this.nodeDatas.splice(idx, 1)
      selNodeData.ui.inLines.slice(0).forEach( line => {
        this.deleteConnect(line)
      })
      selNodeData.ui.outLines.slice(0).forEach( line => {
        this.deleteConnect(line)
      })
      this.$store.dispatch('graph/update')
    },
    deleteSelected() {
      this.$store.state.graph.selNodeDatas.forEach(selNodeData => {
        this.deleteNode(selNodeData)
      })
      this.$store.state.graph.selConnectionDatas.forEach(selConnectionData => {
        this.deleteConnect(selConnectionData)
      })
    },
    keydownHandler(e) {
      let targetTagName = e.target.tagName.toLowerCase()
      if((targetTagName === 'textarea') || (targetTagName === 'input') || !this.$store.getters["graph/editable"]) {
        return;
      }
      if((e.code === 'Delete') || (e.code === 'Backspace')) {
        this.deleteSelected()
      }else if(e.ctrlKey && (e.code === 'KeyC')) {
        this.copySelected()
      }else if(e.ctrlKey && (e.code === 'KeyV')) {
        this.pasteSelected()
      }
    },
    contextMenuClickHandler(menuItem) {
      if(menuItem.key === 'delete') {
        this.deleteSelected()
        this.$store.dispatch('graph/update')
      }else if(menuItem.key === 'copy') {
        this.copySelected()
      }else if(menuItem.key === 'paste') {
        this.pasteSelected()
      }else if(menuItem.key === 'data-view') {
        this.nodeDataVisible = true
      }
    },
    copySelected() {
      // if(GraphUtils.hasStreamOutNode(this.$store.state.graph.selNodeDatas)) {
      //   message.error('输出组件不能被复制')
      //   return;
      // }
      this.copiedEntities = GraphUtils.copyEntities(this.$store.state.graph.selNodeDatas, this.$store.state.graph.selConnectionDatas)
    },
    pasteSelected() {
      if(!this.copiedEntities) {
        return;
      }

      let { nodes,  connectors } = GraphUtils.pasteEntities(this.copiedEntities)

      if(nodes.length > 0) {
        this.$store.state.graph.selNodeDatas.forEach(node => node.ui.selected = false)
        this.$store.commit('graph/nodeDatas', this.nodeDatas.concat(nodes))
        this.$store.commit('graph/selNodeDatas', nodes)
        this.$store.state.graph.selNodeDatas.forEach(node => node.ui.selected = true)
      }
      if(connectors.length > 0) {
        this.$store.state.graph.selConnectionDatas.forEach(connect => connect.selected = false)
        this.$store.commit('graph/connectDatas', this.connectDatas.concat(connectors))
        this.$store.commit('graph/selConnectionDatas', connectors)
        this.$store.state.graph.selConnectionDatas.forEach(connect => connect.selected = true)
      }
      this.$store.dispatch('graph/update')
    },
    getSvgPos(clientX, clientY) {
      let svgPos = GraphUtils.containerCoordinate2Svg(
          clientX - this.$store.state.graph.graphBounding.left, 
          clientY - this.$store.state.graph.graphBounding.top, 
          this.graphTransX, 
          this.graphTransY, 
          this.graphScale
        );
      let svgX = svgPos.x - GraphUtils.NodeWidth * 0.5,
        svgY = svgPos.y - GraphUtils.NodeHeight * 0.5;
      return {x: svgX, y: svgY}
    }
  }
}
</script>

<style>

</style>
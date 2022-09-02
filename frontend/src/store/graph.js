import * as GraphUtils  from '../utils/graph'
import { debounce, deepCopy } from '../utils'
import { requiredValidate, rangeValidate } from '../utils/validate'
import { updateGraph as updateGraphService } from '../service'


const GraphScales = [0.2, 0.4, 0.6, 0.8, 1.0, 1.5, 2.0, 2.5, 3.0, 3.5, 4.0]

const updateGraphImmediate = function(state) {
  if(state.graphStatus === 1) {
    return;
  }
  // 只在编辑的时候同步graph
  state.graphMode = 2
  updateGraphService(GraphUtils.toGraphRawData(state)).then(() => {
    state.graphMode = 3
  }).catch(err => {
    console.error('graph update error', err)
    state.graphMode = 4
  })
}
const updateGraph = debounce(updateGraphImmediate, 1500)

export default {
  namespaced: true,
  state: {
    componentRawData: [],
    nodeDatas: [],
    connectDatas: [],
    graphScales: GraphScales,
    graphScale: GraphScales[4],
    graphScaleIndex: 4,
    graphTransX: 0,
    graphTransY: 0,
    graphBounding: {left: 0, right: 0, top: 0, bottom: 0},
    selNodeDatas: [],
    selConnectionDatas: [],
    graphStatus: 0,     // 0: 编辑  1: 部署
    graphMode: 0,       // 0: 空闲 1: 编辑中 2: 保存中 3：保存成功 4：保存失败
    graphLoading: false,
    processStatus: 0,    // 当前graph的运行状态 0:已停止; 1:运行中
    processMode: 0,     // 当前graph的运行模式 0:停止运行; 1:全部运行; 2:部分运行
  },
  getters: {
    editable(state) {
      return (state.graphStatus === 0) && (state.processStatus === 0)
    }
  },
  mutations: {
    componentRawData(state, val) {
      state.componentRawData = val
    },
    graphBounding(state, val) {
      state.graphBounding = val
    },
    nodeDatas(state, val) {
      state.nodeDatas = val
    },
    connectDatas(state, val) {
      state.connectDatas = val
    },
    selNodeDatas(state, val) {
      state.selNodeDatas = val
    },
    selConnectionDatas(state, val) {
      state.selConnectionDatas = val
    },
    graphStatus(state, val) {
      state.graphStatus = val
    },
    graphLoading(state, val) {
      state.graphLoading = val
    },
    processStatus(state, val) {
      state.processStatus = val
    },
    processMode(state, val) {
      state.processMode = val
    },
    clean(state) {
      state.nodeDatas = []
      state.connectDatas = []
    }
  },
  actions: {
    validateGraph({ state, commit }) {
      return new Promise((resolve, reject) => {
        if(state.nodeDatas.length < 1) {
          reject('部署失败，图中没有节点')
        }else {
          // let errMsg = []
          // for(let i = 0; i < state.nodeDatas.length; i++) {
          //   let nodeData = state.nodeDatas[i]
          //   for(let j = 0; j < nodeData.parameters.length; j++) {
          //     let param = nodeData.parameters[j]
          //     if(!requiredValidate(param.value, param, nodeData.parameters) || !rangeValidate(param.value, param)) {
          //       errMsg.push(`"${nodeData.name}"节点中的参数不满足要求`)
          //       break
          //     }
          //   }
          // }
          if(state.nodeDatas.every(nodeData => nodeData.ui.valid)) {
            resolve()
          }else {
            reject('部署失败，图中有节点的参数不满足要求')
          }
        }
      })
    },
    update({state}) {
      updateGraph(state)
    },
    generateGraph({state}, graphRawData) {
      let componentsMap = {}
      state.componentRawData.forEach(comp => {
        componentsMap[comp.key] = comp
      })

      if(graphRawData.scale != null) {
        state.graphScale = graphRawData.scale
        state.graphScaleIndex = state.graphScales.indexOf(graphRawData.scale)
      }
      if(graphRawData.x != null) {
        state.graphTransX = graphRawData.x
      }
      if(graphRawData.y != null) {
        state.graphTransY = graphRawData.y
      }

      let nodesMap = {}

      let nodeDatas = [];
      let graphNodes = graphRawData.nodes || [];
      for(let i = 0; i < graphNodes.length; i++) {
        if(componentsMap[graphNodes[i].key] == null) {
          continue;
        }
        let nodeData = deepCopy(componentsMap[graphNodes[i].key])
        let nodeDataParams = nodeData.parameters || []
        let graphNodeParams = graphNodes[i].parameters || []
        
        graphNodeParams.forEach(graphNodeParam => {
          let tgtNodeParam = nodeDataParams.find(nodeDataParam => nodeDataParam.key === graphNodeParam.key)
          if(tgtNodeParam) {
            tgtNodeParam.value = graphNodeParam.value
          }
        })
        let node = GraphUtils.generateNode(nodeData, { x: graphNodes[i].x, y:graphNodes[i].y, uuid: graphNodes[i].uuid })
        nodesMap[node.uuid] = node
        nodeDatas.push(node)
      }

      let connectDatas = []
      let graphConnectors = graphRawData.connectors || [];
      for(let i = 0; i < graphConnectors.length; i++) {
        let graphConnector = graphConnectors[i]
        let outNode = nodesMap[graphConnector.src.uuid],
          inNode = nodesMap[graphConnector.tgt.uuid];
        if((outNode == null) || (inNode == null)) {
          continue
        }
        if((inNode.ports.out == null) || (outNode.ports.in == null)) {
          continue
        }
        let inPort = inNode.ports.in.find(port => port.id === graphConnector.tgt.port),
          outPort = outNode.ports.out.find(port => port.id === graphConnector.src.port);
        
        if(inPort && outPort) {
          inPort.matched = true
          outPort.matched = true
          connectDatas.push(GraphUtils.generateConnect(inPort,inNode,outPort,outNode))
        }
      }
      
      state.nodeDatas = nodeDatas
      state.connectDatas = connectDatas

      nodeDatas = null
      connectDatas = null
      componentsMap = null
      nodesMap = null
    },
    validateGraphBeforeDrop({state}, droppingNodes) {
      return new Promise( (resolve, reject) => {
        if(!Array.isArray(droppingNodes)) {
          droppingNodes = [droppingNodes]
        }
        let valid = true
        let droppingNode = null
        for(let i = 0, len = droppingNodes.length; i < len; i++) {
          droppingNode = droppingNodes[i]
          if(droppingNode.key.startsWith('StreamOut')) {
            if(state.nodeDatas.some(nodeData => nodeData.key === droppingNode.key)) {
              valid = false
              break
            }
          }
        }
        if(valid) {
          resolve()
        }else {
          reject(`${droppingNode.name}组件已存在`)
        }
      })
    }
  }
};


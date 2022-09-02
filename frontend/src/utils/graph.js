import { markRaw } from 'vue'
import { v4 as uuidv4 } from 'uuid';
import { deepCopy } from './index';
import { requiredValidate, rangeValidate,
  requiredMsg, outRangeMsg} from './validate';

export const NodeWidth = 230
export const NodeHeight = 44
export const NodeTextWidth = 156

export function contains(x, y, rx, ry, rw, rh) {
  return (x > rx) && ( x < rx + rw) && (y > ry) && (y < ty + rh)
}

export function contains2(x, y, left, right, top, bottom) {
  return (x > left) && (x < right) && (y > top) && (y < bottom)
}

// 两种坐标系，containerCoordinate 和 svgCoordinate
export function scaleBypoint(px, py, translateX, translateY, currentScale, targetScale) {
  let pointerInSvg = containerCoordinate2Svg(px, py, translateX, translateY, currentScale);
  return {
    scale: targetScale,
    translateX: -pointerInSvg.x * targetScale + px,
    translateY: -pointerInSvg.y * targetScale + py
  }
}

export function containerCoordinate2Svg(x, y, translateX, translateY, scale) {
  return {
    x: (x - translateX) / scale,
    y: (y - translateY) / scale
  }
}

export function getNodeUUID() {
  return uuidv4().split('-').join('')
}

export function textEllipsis(text) {
  const svgNS = "http://www.w3.org/2000/svg";
  // let txtNode = document.createTextNode(text);
  let textNode = document.createElementNS(svgNS, "text");
  textNode.textContent = text
  textNode.setAttributeNS(null,"x",100);
  textNode.setAttributeNS(null,"y",100);
  textNode.setAttributeNS(null,"fill","black");
  let svg = document.getElementById("ml-graph");                                  
  svg.appendChild(textNode);

  while (textNode.getComputedTextLength() > NodeTextWidth) {
    text = text.slice(0,-1);
    textNode.textContent = text + "...";
  }
  
  if(text !== textNode.textContent) {
    text = text + "..."
  }

  svg.removeChild(textNode);

  return text;
}


export function toGraphRawData(graphState) {
  
  let nodes = graphState.nodeDatas.map(nd => {
    return {
      uuid: nd.uuid,
      type: nd.type,
      key: nd.key,
      name: nd.name,
      x: nd.x,
      y: nd.y,
      parameters: nd.parameters.map(param => {
        return {
          key: param.key,
          value: param.value
        }
      })
    }
  })
  let connectors = graphState.connectDatas.map( cd => {
    return {
      src: {
        uuid: cd.outNode.uuid,
        name: cd.outNode.name,
        port: cd.outPort.id
      },
      tgt: {
        uuid: cd.inNode.uuid,
        name: cd.inNode.name,
        port: cd.inPort.id
      }
    }
  })
  return {
    htime: Date.now(),
    scale: graphState.graphScale,
    x: graphState.graphTransX,
    y: graphState.graphTransY,
    nodes,
    connectors
  }
}

export function generateConnect(inPort,inNode,outPort,outNode) {
  inPort.matched = true
  outPort.matched = true
  let connector = {
    uuid: getNodeUUID(),
    selected: false,
    inPort,
    inNode,
    outPort,
    outNode
  }
  nodeAddConnector(inNode.ui.inLines, 'in', connector)
  nodeAddConnector(outNode.ui.outLines, 'out', connector)
  return connector
}

export function getNodeValidation(parameters) {
  return parameters.every(param => param.valid)
}

export function generateNode(nodeData, {x, y, uuid}) {
  nodeData.displayName = textEllipsis(nodeData.name)

  let nodeValid = true
  nodeData.parameters.forEach((param, idx, params) => {
    if(param.value == null) {
      param.value = null
      if(param.default != null) {
        param.value = deepCopy(param.default)
      }
    }
    param.valid = true
    param.errMsg = ''
    if(!requiredValidate(param.value, param, params)) {
      param.valid = false
      param.errMsg = requiredMsg(param)
    }
    if(param.valid) {
      if(!rangeValidate(param.value, param, params)) {
        param.valid = false
        param.errMsg = outRangeMsg(param)
      }
    }
    nodeValid = nodeValid && param.valid
  })

  if(nodeData.ports == null) {
    nodeData.ports = { in: [], out: []}
  }

  nodeData.ports.in = nodeData.ports.in || []
  nodeData.ports.out = nodeData.ports.out || []

  if(nodeData.key === 'ExecutePythonScript' || nodeData.key === 'DataProcess') {
    for(let i = 0, len = nodeData.parameters.length; i < len; i++) {
      let param = nodeData.parameters[i]
      if((param.key === 'inPorts') && (Array.isArray(param.value))) {
        nodeData.ports.in = (deepCopy(param.value))
      }
      if((param.key === 'outPorts') && (Array.isArray(param.value))) {
        nodeData.ports.out = (deepCopy(param.value))
      }
    }
  }
  
  nodeData.ports.in.forEach((port, idx) => {
    port.x = NodeWidth / (nodeData.ports.in.length + 1) * (idx + 1)
    port.y = 0
    port.matchingType = null
    port.matched = false
  })

  nodeData.ports.out.forEach((port, idx) => {
    port.x = NodeWidth / (nodeData.ports.out.length + 1) * (idx + 1)
    port.y = NodeHeight
    port.matchingType = null
    port.matched = false
  })

  Object.assign(
    nodeData, 
    {
      uuid: uuid ? uuid : getNodeUUID(),
      x: x,
      y: y
    }, 
    {
      ui: {
        valid: nodeValid,
        selected: false,
        matchType: null,
        inLines: markRaw([]),
        outLines: markRaw([])
      }
    }
  )
  updatePortsPosition(nodeData);
  return nodeData
}

export function deleteConnect(connectData, connectDatas) {
  let idx = connectDatas.findIndex(
    c => ((c.inNode.uuid === connectData.inNode.uuid) && (c.outNode.uuid === connectData.outNode.uuid) && 
      (c.inPort.id === connectData.inPort.id) && (c.outPort.id === connectData.outPort.id))
  )
  if(idx < 0) {
    return
  }

  connectDatas.splice(idx, 1)
  
  nodeDeleteConnector(connectData.inNode.ui.inLines, 'in', connectData)
  nodeDeleteConnector(connectData.outNode.ui.outLines, 'out', connectData)

  // update port state
  if(nodeGetPortConnectorsByPort(connectData.inNode.ui.inLines, 'in', connectData.inNode, connectData.inPort).length > 0) {
    connectData.inPort.matched = true
  }else {
    connectData.inPort.matched = false
  }
  if(nodeGetPortConnectorsByPort(connectData.outNode.ui.outLines, 'out', connectData.outNode, connectData.outPort).length > 0) {
    connectData.outPort.matched = true
  }else {
    connectData.outPort.matched = false
  }
}

export function addPort(type, nodeData, { name, id }) {
  let port = {
    id,
    name,
    matchingType: null,
    matched: false
  }
  let ports = []
  if(type === 'in') {
    nodeData.ports.in.push(port)
    ports = nodeData.ports.in
  }else if(type === 'out') {
    nodeData.ports.out.push(port)
    ports = nodeData.ports.out
  }
  ports.forEach((port, idx, arr) => {
    port.x = NodeWidth / (arr.length + 1) * (idx + 1)
    port.y = type === 'in' ? 0 : NodeHeight,
    port.matchingType = null
    port.matched = false
  })
  updatePortsPosition(nodeData, type)
}

export function deleteLastPort(type, nodeData, connectDatas) {
  let deletingPort = null, nodeLines = [], ports = [];
  if(type === 'in') {
    deletingPort = nodeData.ports.in[nodeData.ports.in.length - 1]
    nodeLines = nodeData.ui.inLines
    ports = nodeData.ports.in
  }else if(type === 'out') {
    deletingPort = nodeData.ports.out[nodeData.ports.out.length - 1]
    nodeLines = nodeData.ui.outLines
    ports = nodeData.ports.out
  }
  if(deletingPort && deletingPort.matched) {
    let portConnectors = nodeGetPortConnectorsByPort(nodeLines, type, nodeData, deletingPort)
    portConnectors.forEach(connector => {
      deleteConnect(connector, connectDatas)
    })
  }
  ports.splice(ports.length - 1, 1);
  ports.forEach((port, idx, arr) => {
    port.x = NodeWidth / (arr.length + 1) * (idx + 1)
  })
  updatePortsPosition(nodeData, type)
}


export function updatePortsPosition(nodeData, type='all') {
  const update = (port, idx, ports) => {
    port.svgX = nodeData.x + port.x
    port.svgY = nodeData.y + port.y
    port.svgBTop = nodeData.y
    port.svgBBottom = nodeData.y + NodeHeight
    updatePortPosition(idx, ports, nodeData)
  }
  if(type === 'in') {
    nodeData.ports.in.forEach(update)
  }else if(type === 'out') {
    nodeData.ports.out.forEach(update)
  }else {
    nodeData.ports.in.forEach(update)
    nodeData.ports.out.forEach(update)
  }
}

export function nodeAddConnector(nodeLines, portType, connector) {
  if(portType === 'in') {
    let idx = nodeLines.findIndex(line => line.uuid === connector.uuid)
    if(idx === -1) {
      nodeLines.push(connector)
    }
  }else if(portType === 'out') {
    let idx = nodeLines.findIndex(line => line.uuid === connector.uuid)
    if(idx === -1) {
      nodeLines.push(connector)
    }
  }
}

export function nodeDeleteConnector(nodeLines, portType, connector) {
  if(portType === 'in') {
    let idx = nodeLines.findIndex(line => line.uuid === connector.uuid)
    if(idx > -1) {
      nodeLines.splice(idx, 1)  
    }
  }else if(portType === 'out') {
    let idx = nodeLines.findIndex(line => line.uuid === connector.uuid)
    if(idx > -1) {
      nodeLines.splice(idx, 1)  
    }
  }
}

export function nodeGetPortConnectorsByPort(nodeLines, portType, connectorNode, connectorPort) {
  let res = []
  if(portType === 'in') {
    res = nodeLines.filter(line => (line.inNode.uuid === connectorNode.uuid) && (line.inPort.id === connectorPort.id))
  }else if(portType === 'out') {
    res = nodeLines.filter(line => (line.outNode.uuid === connectorNode.uuid) && (line.outPort.id === connectorPort.id))
  }
  return res;
}

function updatePortPosition(idx, ports, nodeData) {
  if(ports.length === 1) {
    ports[idx].svgBLeft = nodeData.x
    ports[idx].svgBRight = nodeData.x + NodeWidth
  }else {
    if(idx === 0) {
      ports[idx].svgBLeft = nodeData.x
      ports[idx].svgBRight = nodeData.x + NodeWidth / (2 * ports.length + 2) * 3
    }else if(idx === ports.length - 1) {
      ports[idx].svgBRight = nodeData.x + NodeWidth
      ports[idx].svgBLeft = nodeData.x + NodeWidth - NodeWidth / (2 * ports.length + 2) * 3
    }else {
      ports[idx].svgBLeft = nodeData.x + NodeWidth / (2 * ports.length + 2) * (3 + idx * 2 - 2)
      ports[idx].svgBRight = nodeData.x + NodeWidth / (2 * ports.length + 2) * (3 + idx * 2)
    }
  }
}

export function validateGraphBeforeDrop(droppingNodes, nodeDatas) {
  // if(!Array.isArray(droppingNodes)) {
  //   droppingNodes = [droppingNodes]
  // }
  // let valid = true
  // let droppingNode = null
  // for(let i = 0, len = droppingNodes.length; i < len; i++) {
  //   droppingNode = droppingNodes[i]
  //   if(droppingNode.key.startsWith('StreamOut')) {
  //     if(nodeDatas.some(nodeData => nodeData.key === droppingNode.key)) {
  //       valid = false
  //       break
  //     }
  //   }
  // }
  // if(valid) {
  //   return {error: null}
  // }else {
  //   return {
  //     error: new Error(`${droppingNode.name}组件已存在`)
  //   }
  // }
  return {error: null}
}

export function hasStreamOutNode(nodeDatas) {
  return nodeDatas.some(nodeData => nodeData.key.startsWith('StreamOut'))
}


export function copyEntities(selectedNodes, selectedConnectors) {
  if(selectedNodes.length < 1) {
    return []
  }
  let selectedNodesMap = {}
  for(let i = 0, len = selectedNodes.length; i < len; i++) {
    selectedNodesMap[selectedNodes[i].uuid] = selectedNodes[i]
  }
  let copiedNodes = [], copiedConnectors = [];
  let copiedNodesMap = {};
  for(let i = 0, len = selectedConnectors.length; i < len; i++) {
    // 遍历选中的连线
    let connector = selectedConnectors[i]
    if(selectedNodesMap[connector.inNode.uuid] && selectedNodesMap[connector.outNode.uuid]) {
      let copiedInNode = null;
      let copiedOutNode = null;
      if(copiedNodesMap[connector.inNode.uuid]) {
        copiedInNode = copiedNodesMap[connector.inNode.uuid]
      }
      if(copiedNodesMap[connector.outNode.uuid]) {
        copiedOutNode = copiedNodesMap[connector.outNode.uuid]
      }
      if(!copiedInNode) {
        copiedInNode = copyNode(connector.inNode);
        copiedNodes.push(copiedInNode)
        copiedNodesMap[connector.inNode.uuid] = copiedInNode
      }
      if(!copiedOutNode) {
        copiedOutNode = copyNode(connector.outNode);
        copiedNodes.push(copiedOutNode)
        copiedNodesMap[connector.outNode.uuid] = copiedOutNode
      }
      let copiedInPort = getPortById(copiedInNode.ports.in, connector.inPort.id);
      let copiedOutPort = getPortById(copiedOutNode.ports.out, connector.outPort.id);
      copiedConnectors.push(generateConnect(copiedInPort, copiedInNode, copiedOutPort, copiedOutNode));
    }
  }
  for(let i = 0, len = selectedNodes.length; i < len; i++) {
    if(copiedNodesMap[selectedNodes[i].uuid]) {
      continue;
    }
    let copiedNode = copyNode(selectedNodes[i]);
    copiedNodesMap[selectedNodes[i].uuid] = copiedNode;
    copiedNodes.push(copiedNode);
  }

  return { nodes: copiedNodes, connectors: copiedConnectors }
}

export function getPortById(ports, portId) {
  return ports.find(port => port.id === portId)
}

export function copyNode(copyingNode) {
  let inLines = copyingNode.ui.inLines
  let outLines = copyingNode.ui.outLines
  copyingNode.ui.inLines = []
  copyingNode.ui.outLines = []
  
  let copiedNode = deepCopy(copyingNode)

  copyingNode.ui.inLines = inLines
  copyingNode.ui.outLines = outLines

  copiedNode.uuid = getNodeUUID()
  copiedNode.ui.selected = false

  copiedNode.ports.in.forEach(port => port.matched = false)
  copiedNode.ports.out.forEach(port => port.matched = false)

  return copiedNode;
}

export function pasteEntities({ nodes: copiedNodes, connectors: copiedConnectors }) {
  let { nodes, connectors } = copyEntities(copiedNodes, copiedConnectors)
  for(let i = 0, len = nodes.length; i < len; i++) {
    let node = nodes[i]
    node.x += 50
    node.y += 20
    updatePortsPosition(node);
  }
  return { nodes, connectors }
}
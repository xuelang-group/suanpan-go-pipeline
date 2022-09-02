import SocketIOClient from './sio'
// import graphData  from './mock/graph.js'

export function getGraph() {
  return SocketIOClient.request({event: 'graph.get'})
}

export function updateGraph(data) {
  return SocketIOClient.request({event: 'graph.update', data: data})
}

export function getGraphStatus() {
  return SocketIOClient.request({event: 'graph.status.get'})
}

export function updateGraphStatus(data) {
  return SocketIOClient.request({event: 'graph.status.set', data: data})
}
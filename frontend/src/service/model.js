import SocketIOClient from './sio'

export function getModelList() {
  return SocketIOClient.request({event: 'graph.model.list'})
}

export function publishModel(data) {
  return SocketIOClient.request({event: 'graph.model.publish', data: data})
}

export function deleteModel(data) {
  return SocketIOClient.request({event: 'graph.model.remove', data: data})
}

export function updateModel(data) {
  return SocketIOClient.request({event: 'graph.model.update', data: data})
}

export function selectModel(data) {
  return SocketIOClient.request({event: 'graph.model.select', data: data})
}

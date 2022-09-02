import SocketIOClient from './sio'

export function runProcess(graphData) {
  return SocketIOClient.emit({event: 'process.run', data: {
    mode: 1,
    graph: graphData
  }})
}

export function stopProcess() {
  return SocketIOClient.emit({event: 'process.run', data: {
    mode: 0
  }})
}

export function getProcessStatus() {
  return SocketIOClient.request({event: 'process.status.get'})
}

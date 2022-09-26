import SocketIOClient from './sio'

export function runProcess() {
  return SocketIOClient.emit({event: 'process.run'})
}

export function stopProcess() {
  return SocketIOClient.emit({event: 'process.stop'})
}

export function getProcessStatus() {
  return SocketIOClient.request({event: 'process.status.get'})
}

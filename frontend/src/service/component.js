import SocketIOClient from './sio'
// import fakeData from './mock/componentList'

export function getComponentList() {
  return SocketIOClient.request({event: 'algos.get'})
  // return Promise.resolve(fakeData)
}
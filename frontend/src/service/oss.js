import SocketIOClient from './sio' 
import OSSService from './ossService/index'

export let ossService = null

export function ossServiceInit(type) {
  if(type === 'oss') {
    ossService = OSSService({type: 'oss', client: true})
  }else {
    ossService = OSSService({type: 'minio', client: true})
  }
}

export function getStorageInfo() {
  return SocketIOClient.request({event: 'storage.info'})
}
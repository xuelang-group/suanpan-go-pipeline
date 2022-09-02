import { createStore } from 'vuex'
import graph from './graph'

export default createStore({
  modules: {
    graph
  },
  state: {
    initLoading: true,
    storageNodePath: null,
  },
  mutations: {
    initLoading(state, val) {
      state.initLoading = val
    },
    storageNodePath(state, val) {
      state.storageNodePath = val
    },
  },
  actions: {
  }
})

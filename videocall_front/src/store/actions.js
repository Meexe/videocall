export default {
  openWebsocket ({ commit }) {
    commit('openWebsocket')
  },
  setWebsocketMessageHandler ({ commit }, handler) {
    commit('setWebsocketMessageHandler', handler)
  },
  closeWebsocket ({ commit }) {
    commit('closeWebsocket')
  }
}
export default {
  openWebsocket (state) {
    state.ws = new WebSocket("ws://127.0.0.1:8000/api/ws/online")
    state.ws.onopen = function(evt) {
      console.log("OPEN " + evt.data)
    }
    state.ws.onclose = function(evt) {
      console.log("CLOSE " + evt.data)
      state.ws = null
    }
    state.ws.onerror = function(evt) {
      console.log("ERROR: " + evt.data)
      state.ws = null
    }
  },
  setWebsocketMessageHandler (state, handler) {
    state.ws.onmessage = handler
  },
  closeWebsocket (state) {
    state.ws.close()
  }
}
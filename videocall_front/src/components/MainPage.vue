<template>
  <div id="videocall">
    <h3>Videocalls</h3>
    <h4>Welcome, {{ name }}</h4>
    <h4>Users online:</h4>
    <div v-for="user in users" :key="user.id">
      <span>
        {{ user }}
        <button @click="call(user)">Call</button>
      </span>
    </div>
  </div>
</template>

<script>
export default {
  name: 'MainPage',
  data () {
      return {
          name: null,
          ws: null,
          users: null,
      }
  },
  methods: {
    call (user) {
      this.ws.send(user)
    }
  },
  mounted () {
    let user = this.$store.getters.user
    if (user === null) {
      this.$router.push('/login')
      return
    }
    let vm = this
    this.name = user.nickname
    this.ws = new WebSocket("ws://127.0.0.1:8000/api/ws/online")
    this.ws.onopen = function(evt) {
      console.log("OPEN " + evt.data)
    }
    this.ws.onclose = function(evt) {
      console.log("CLOSE " + evt.data)
      vm.ws = null
    }
    this.ws.onmessage = function(evt) {
      let msg = JSON.parse(evt.data)
      console.log(msg)
      if (msg.type === 'users') vm.users = msg.payload
      else if (msg.type === 'call') alert(msg.payload)
    }
    this.ws.onerror = function(evt) {
      console.log("ERROR: " + evt.data)
      vm.ws = null
    }
  },
  beforeDestroy () {
    if (this.ws !== null) {
      this.ws.close()
    }
  }
}
</script>

<style scoped>
</style>

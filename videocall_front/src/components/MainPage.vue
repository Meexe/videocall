<template>
  <div id="main">
    <h3>Videocalls</h3>
    <h4>Welcome, {{ localUser }}</h4>
    <h4>Users online:</h4>
    <div v-for="user in users" :key="user.id">
      <span>
        {{ user }}
        <button @click="call(user)">Call</button>
      </span>
    </div>
    <div v-if="showModal">
      {{ textModal }}
      <div v-if="isCall">
        <button @click="acceptCall">
          Accept
        </button>
        <button @click="declineCall">
          Decline
        </button>
      </div>
      <div v-else>
        <button @click="closeModal">
          OK
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'MainPage',
  data () {
    return {
      localUser: null,
      remoteUser: null,
      users: null,
      showModal: false,
      isCall: false,
      textModal: null,
    }
  },
  computed: {
    ws () {
      return this.$store.getters.ws
    }
  },
  methods: {
    setStandardHandler () {
      let vm = this
      vm.$store.dispatch('setWebsocketMessageHandler', function(evt) {
        let msg = JSON.parse(evt.data)
        console.log(msg)
        if (msg.type === 'users') vm.users = msg.payload
        else if (msg.type === 'call' && msg.payload === 'propose') {
          vm.remoteUser = msg.source
          vm.openModal('user ' + msg.source + ' is calling you', true)
        }
      })
    },
    call (user) {
      let vm = this
      this.remoteUser = user
      let req = {
        type: 'call',
        source: this.localUser,
        destination: this.remoteUser,
        payload: 'propose'
      }
      this.ws.send(JSON.stringify(req))
      let timer = setTimeout(function() {
        console.log('call timed out')
        vm.openModal('call timed out', false)
        vm.remoteUser = null
        vm.setStandardHandler()
      }, 15000)
      this.$store.dispatch('setWebsocketMessageHandler', function(evt) {
        let msg = JSON.parse(evt.data)
        console.log(msg.type, msg.payload)
        if (msg.type === 'call' && msg.payload === 'accept') {
          console.log('accepted')
          clearTimeout(timer)
          vm.$router.push({ path: `/videocall/${vm.localUser}&${vm.remoteUser}` })
        } else if (msg.type === 'call' && msg.payload === 'decline') {
          console.log('declined')
          clearTimeout(timer)
          vm.openModal('user declined your call', false)
          vm.setStandardHandler()
          vm.remoteUser = null
        }
      })
    },
    openModal (msg, isCall) {
      this.isCall = isCall
      this.showModal = true
      this.textModal = msg
    },
    closeModal () {
      this.isCall = false
      this.showModal = false
      this.textModal = null
    },
    acceptCall () {
      let msg = {
        type: 'call',
        source: this.localUser,
        destination: this.remoteUser,
        payload: 'accept'
      }
      this.ws.send(JSON.stringify(msg))
      this.$router.push({ path: `/videocall/${this.localUser}&${this.remoteUser}` })
    },
    declineCall () {
      let msg = {
        type: 'call',
        source: this.localUser,
        destination: this.remoteUser,
        payload: 'decline'
      }
      this.ws.send(JSON.stringify(msg))
      this.remoteUser = null
      this.closeModal()
    }
  },
  mounted () {
    this.localUser = this.$route.params.user
    this.$store.dispatch('openWebsocket')
    this.setStandardHandler()
  }
}
</script>

<style scoped>
</style>

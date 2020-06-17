<template>
  <div id="videocall">
    <div id="videos">
      <video id="localVideo" muted autoplay playsinline></video>
      <video id="remoteVideo" autoplay playsinline></video>
    </div>
  </div>
</template>

<script>
export default {
  name: 'VideocallPage',
  data () {
    return {
      localUser: null,
      remoteUser: null,
      pc: null
    }
  },
  computed: {
    ws () {
      return this.$store.getters.ws
    }
  },
  mounted () {
    this.localUser = this.$route.params.localUser
    this.remoteUser = this.$route.params.remoteUser
    this.pc = new RTCPeerConnection({
      iceServers: [
        {
          urls: 'stun:stun.l.google.com:19302'
        }
      ]
    })

  }
}
</script>

<style scoped>
</style>

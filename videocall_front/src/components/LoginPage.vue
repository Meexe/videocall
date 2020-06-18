<template>
  <div class="login">
    <h4>Login</h4>
    <form>
      <label for="nickname">Nickname</label>
      <div>
        <input id="nickname" type="nickname" v-model="nickname" required autofocus>
      </div>
      <div>
        <label for="password">Password</label>
        <div>
          <input id="password" type="password" v-model="password" required>
        </div>
      </div>
      <div>
        <button type="submit" @click="handleSubmit">
          Login
        </button>
      </div>
      <br><br><br><br>
      <router-link to="/signup">Signup</router-link>
    </form>
  </div>
</template>

<script>
const axios = require('axios');
export default {
  name: 'LoginPage',
  data () {
    return {
      nickname : "",
      password : ""
    }
  },
  methods: {
    handleSubmit(e) {
      e.preventDefault()
      if (this.password.length > 5) {
        axios.post('https://max.pak.digital/api/user/login', {
          nickname: this.nickname,
          password: this.password
        })
        .then(response => {
          console.log(response.data)
          if (response.data.status) {
            document.cookie = 'jwt=' + response.data.user.token
            this.$router.push({ path: `/user/${response.data.user.nickname}` })
          } else {
            console.log(response.message)
          }
        })
        .catch(function (error) {
          console.error(error.response);
        });
      }
    }
  }
}
</script>

<style scoped>
</style>
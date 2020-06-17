import Vue from 'vue'
import Vuex from 'vuex'
import VueRouter from 'vue-router'
import store from './store'
import App from './components/App.vue'
import MainPage from './components/MainPage.vue'
import LoginPage from './components/LoginPage.vue'
import SignupPage from './components/SignupPage.vue'
import VideocallPage from './components/VideocallPage.vue'

Vue.use(Vuex)
Vue.use(VueRouter)
Vue.config.productionTip = false

const routes = [
  { path: '/user/:user', component: MainPage },
  { path: '/', component: LoginPage },
  { path: '/signup', component: SignupPage },
  { path: '/videocall/:localUser&:remoteUser', component: VideocallPage }
]

const router = new VueRouter({
  routes
})

new Vue({
  store,
  router,
  render: h => h(App)
}).$mount('#app')

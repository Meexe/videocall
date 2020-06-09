import Vue from 'vue'
import Vuex from 'vuex'
import VueRouter from 'vue-router'
import store from './store'
import App from './components/App.vue'
import MainPage from './components/MainPage.vue'
import LoginPage from './components/LoginPage.vue'
import SignupPage from './components/SignupPage.vue'

Vue.use(Vuex)
Vue.use(VueRouter)
Vue.config.productionTip = false

const routes = [
  { path: '/', component: MainPage },
  { path: '/login', component: LoginPage },
  { path: '/signup', component: SignupPage },
  // { path: '/videocall', component: VideocallPage }
]

const router = new VueRouter({
  routes
})

new Vue({
  store,
  router,
  render: h => h(App)
}).$mount('#app')

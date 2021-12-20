
import { createApp } from 'vue'
import App from 'src/App.vue'
import '@quasar/extras/material-icons/material-icons.css'
import axios from 'axios'
import Home from './views/Home.vue'
import TournamentView from './views/TournamentView.vue'
//import TournamentList from './components/TournamentList.vue'
import VueKeycloakJs from '@dsb-norge/vue-keycloak-js'
import { KeycloakInstance } from "keycloak-js";
import { VueKeycloakInstance } from "@dsb-norge/vue-keycloak-js/dist/types";
import router from './router.js'
import { Quasar } from 'quasar'

//const app = createApp(App)
import 'quasar/src/css/index.sass'



const axiosApiInstance = axios.create({
  baseURL: 'http://localhost:3000/api',
})

function tokenInterceptor() {
  axiosApiInstance.interceptors.request.use(
    async config => {
      const token = app.config.globalProperties.$keycloak.token
      config.headers = {
        Authorization: `Bearer ` + token
      }
      return config
    },
    error => {
      Promise.reject(error)
    })
}

const app = createApp(App)
app.use(router)
app.use(Quasar, {
  plugins: {}, // import Quasar plugins and add here
})
app.provide('router', router)
app.use(VueKeycloakJs, {
  init: {
    //onLoad: 'check-sso',
    onLoad: 'login-required',
    silentCheckSsoRedirectUri: window.location.origin + "/silent-check-sso.html"
  },
  config: {
    url: 'http://localhost:8080/auth',
    clientId: 'brackets-app',
    realm: 'GoogleBrackets',
    onAuthRefreshSuccess: function () { console.log("AuthRefreshSuccess") }
  },
  onReady(keycloak: KeycloakInstance) {
    console.log('Keycloak ready')
    tokenInterceptor()
    app.provide('keycloak', app.config.globalProperties.$keycloak)
    app.component('Home', Home)
    app.component('TournamentView', TournamentView)
///    app.component('TournamentList', TournamentList)
    app.mount('#app');

  }
})

// install all modules under `modules/`
//Object.values(import.meta.globEager('./plugins/*.ts')).map((plugin) => plugin.install?.(app))

export { axiosApiInstance }
///app.mount('#app')

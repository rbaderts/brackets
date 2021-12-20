import axios from 'axios'

import App from "../App.vue"

const axiosApiInstance = axios.create({
      baseURL: 'http://localhost:3000/api',
})

// Request interceptor for API calls
axiosApiInstance.interceptors.request.use(

  async config => {
//    let = getKeycloak();
//    const token = getToken()
    const token = App.config.globalProperties.$keycloak.token

    console.log("token = "); console.log(token)
    config.headers = {
      Authorization: `Bearer ` + token
    }
    return config
  },
  error => {
    Promise.reject(error)
  },
)

export default axiosApiInstance

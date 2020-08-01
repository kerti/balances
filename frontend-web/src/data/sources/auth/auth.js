import axios from 'axios'
import sources from '..'

export default {
  login: (username, password) => {
    const basic = btoa(`${username}:${password}`)
    return axios.post(`${sources.baseURL}auth/login`, null, {
      headers: {
        Authorization: `Basic ${basic}`,
      },
    })
  },
}

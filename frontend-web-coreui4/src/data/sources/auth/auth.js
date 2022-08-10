import axios from 'axios'
import sources from '..'

const auth = {
  login: (username, password) => {
    const basic = btoa(`${username}:${password}`)
    return axios.post(`${sources.baseURL}auth/login`, null, {
      headers: {
        Authorization: `Basic ${basic}`,
      },
    })
  },
}

export default auth

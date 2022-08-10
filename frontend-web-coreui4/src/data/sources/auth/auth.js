import axios from 'axios'
import sources from '../sources'

const auth = {
  login: (username, password) => {
    const basic = btoa(`${username}:${password}`)
    return axios.post(`${sources.baseURL}/auth/login`, null, {
      headers: {
        Authorization: `Basic ${basic}`,
        Origin: sources.origin,
      },
    })
  },
}

export default auth
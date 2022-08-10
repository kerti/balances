import types from './types'
import auth from '../sources/auth/auth'
import Cookies from 'cookies-js'
import cookieNames from '../cookies'

export function requestLogin(username, password, navigate) {
  return (dispatch) => {
    dispatch(loginLoading())
    auth
      .login(username, password)
      .then((payload) => {
        const loginResponse = payload.data.data

        Cookies.set(cookieNames.auth.token, loginResponse.token, {
          expires: new Date(loginResponse.expiration),
        })

        dispatch(loginSuccess(navigate, payload.data))
        navigate('/')
      })
      .catch((error) => {
        dispatch(loginFailure(error.response.data))
      })
  }
}

export function loadAuthCookies() {
  const payload = {
    data: {
      token: Cookies.get(cookieNames.auth.token),
    },
  }
  return {
    type: types.auth.login.LOADCOOKIES,
    payload: payload,
  }
}

export function loginLoading() {
  return {
    type: types.auth.login.LOADING,
  }
}

export function loginSuccess(navigate, payload) {
  navigate('/')
  return {
    type: types.auth.login.SUCCESS,
    payload: payload,
  }
}

export function loginFailure(error) {
  return {
    type: types.auth.login.FAILURE,
    error: error,
  }
}

export function requestLogout(navigate) {
  navigate('/login')

  Cookies.expire(cookieNames.auth.token)

  return {
    type: types.auth.logout.REQUEST,
  }
}

export function requestAuthCheck(navigate) {
  const token = Cookies.get(cookieNames.auth.token)
  return (dispatch) => {
    // TODO: send refresh token here instead of forcing a logout?
    if (!token) dispatch(requestLogout(navigate))
  }
}

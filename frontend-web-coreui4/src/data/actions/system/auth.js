import { actionTypes } from '../../actions'
import auth from '../../sources/auth/auth'
import Cookies from 'cookies-js'
import cookieNames from '../../cookies'

export function requestLogin(username, password, navigate) {
  return (dispatch) => {
    dispatch(loginLoading())
    auth
      .login(username, password)
      .then((payload) => {
        // TODO: figure out a way to solve the samesite cookie issue
        const loginResponse = payload.data.data

        Cookies.set(cookieNames.auth.token, loginResponse.token, {
          expires: new Date(loginResponse.expiration),
        })

        dispatch(loginSuccess(payload.data))
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
    type: actionTypes.auth.login.LOADCOOKIES,
    payload: payload,
  }
}

export function loginLoading() {
  return {
    type: actionTypes.auth.login.LOADING,
  }
}

export function loginSuccess(payload) {
  return {
    type: actionTypes.auth.login.SUCCESS,
    payload: payload,
  }
}

export function loginFailure(error) {
  return {
    type: actionTypes.auth.login.FAILURE,
    error: error,
  }
}

export function requestLogout(navigate) {
  navigate('/login')

  Cookies.expire(cookieNames.auth.token)

  return {
    type: actionTypes.auth.logout.REQUEST,
  }
}

export function requestAuthCheck(navigate) {
  const token = Cookies.get(cookieNames.auth.token)
  return (dispatch) => {
    // TODO: send refresh token here instead of forcing a logout?
    if (!token) dispatch(requestLogout(navigate))
  }
}

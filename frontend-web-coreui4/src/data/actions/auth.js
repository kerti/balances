import actionTypes from './actionTypes'
import auth from '../sources/auth/auth'
import cookieNames from '../cookies'

export function requestLogin(setCookie, username, password, navigate) {
  return (dispatch) => {
    dispatch(loginLoading())
    auth
      .login(username, password)
      .then((payload) => {
        const loginResponse = payload.data.data

        setCookie(cookieNames.auth.token, loginResponse.token, {
          expires: new Date(loginResponse.expiration),
        })

        setCookie(cookieNames.auth.userId, loginResponse.user.id, {
          expires: new Date(payload.data.data.expiration),
        })

        setCookie(cookieNames.auth.userEmail, loginResponse.user.email, {
          expires: new Date(payload.data.data.expiration),
        })

        setCookie(cookieNames.auth.userProfileName, loginResponse.user.name, {
          expires: new Date(payload.data.data.expiration),
        })

        setCookie(cookieNames.auth.username, loginResponse.user.username, {
          expires: new Date(payload.data.data.expiration),
        })

        dispatch(loginSuccess(navigate, payload.data))
      })
      .catch((error) => {
        dispatch(loginFailure(error.response.data))
      })
  }
}

export function loginLoading() {
  return {
    type: actionTypes.auth.login.LOADING,
  }
}

export function loginSuccess(navigate, payload) {
  navigate('/')
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

export function requestLogout(removeCookie, navigate) {
  navigate('/logout')

  removeCookie(cookieNames.auth.token)
  removeCookie(cookieNames.auth.userId)
  removeCookie(cookieNames.auth.userEmail)
  removeCookie(cookieNames.auth.userProfileName)
  removeCookie(cookieNames.auth.username)

  return {
    type: actionTypes.auth.logout.REQUEST,
  }
}

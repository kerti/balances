import actionTypes from "./actionTypes";
import auth from "../sources/auth/auth";

export function requestLogin(setCookies, username, password) {
  return (dispatch) => {
    dispatch(loginLoading());
    auth
      .login(username, password)
      .then((payload) => {
        setCookies("token", payload.data.data.token, {
          expires: new Date(payload.data.data.expiration),
        });
        dispatch(loginSuccess(payload.data));
      })
      .catch((error) => {
        console.log(error);
        dispatch(loginFailure(error.response.data));
      });
  };
}

export function loginLoading() {
  return {
    type: actionTypes.auth.login.LOADING,
  };
}

export function loginSuccess(payload) {
  return {
    type: actionTypes.auth.login.SUCCESS,
    payload: payload,
  };
}

export function loginFailure(error) {
  return {
    type: actionTypes.auth.login.FAILURE,
    error: error,
  };
}

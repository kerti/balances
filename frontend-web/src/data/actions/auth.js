import actionTypes from "./actionTypes";
import auth from "../sources/auth/auth";

export function requestLogin(username, password) {
  return (dispatch) => {
    dispatch(loginLoading);
    auth
      .login(username, password)
      .then(({ data }) => {
        dispatch(loginSuccess(data));
      })
      .catch((error) => {
        dispatch(loginFailure(error));
      });
  };
}

export function loginLoading() {
  return {
    type: actionTypes.auth.login.LOADING,
  };
}

export function loginSuccess(data) {
  return {
    type: actionTypes.auth.login.SUCCESS,
    data: data,
  };
}

export function loginFailure(error) {
  return {
    type: actionTypes.auth.login.FAILURE,
    error: error,
  };
}

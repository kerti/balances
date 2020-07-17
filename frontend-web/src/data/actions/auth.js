import actionTypes from "./actionTypes";
import auth from "../sources/auth/auth";
import Cookies from "cookies-js";
import cookieNames from "../cookies";

export function requestLogin(username, password, history) {
  return (dispatch) => {
    dispatch(loginLoading());
    auth
      .login(username, password)
      .then((payload) => {
        const loginResponse = payload.data.data;

        Cookies.set(cookieNames.auth.token, loginResponse.token, {
          expires: new Date(loginResponse.expiration),
        });

        dispatch(loginSuccess(history, payload.data));
        history.push("/");
      })
      .catch((error) => {
        dispatch(loginFailure(error.response.data));
      });
  };
}

export function loadAuthCookies() {
  const payload = {
    data: {
      token: Cookies.get(cookieNames.auth.token),
    },
  };
  return {
    type: actionTypes.auth.login.LOADCOOKIES,
    payload: payload,
  };
}

export function loginLoading() {
  return {
    type: actionTypes.auth.login.LOADING,
  };
}

export function loginSuccess(history, payload) {
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

export function requestLogout(history) {
  history.push("/login");

  Cookies.expire(cookieNames.auth.token);

  return {
    type: actionTypes.auth.logout.REQUEST,
  };
}

export function requestAuthCheck(history) {
  const token = Cookies.get(cookieNames.auth.token);
  return (dispatch) => {
    // TODO: send refresh token here instead of forcing a logout?
    if (!token) dispatch(requestLogout(history));
  };
}

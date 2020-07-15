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

        Cookies.set(cookieNames.auth.userId, loginResponse.user.id, {
          expires: new Date(payload.data.data.expiration),
        });

        Cookies.set(cookieNames.auth.userEmail, loginResponse.user.email, {
          expires: new Date(payload.data.data.expiration),
        });

        Cookies.set(cookieNames.auth.userProfileName, loginResponse.user.name, {
          expires: new Date(payload.data.data.expiration),
        });

        Cookies.set(cookieNames.auth.username, loginResponse.user.username, {
          expires: new Date(payload.data.data.expiration),
        });

        dispatch(loginSuccess(history, payload.data));
      })
      .catch((error) => {
        dispatch(loginFailure(error.response.data));
      });
  };
}

export function loginLoading() {
  return {
    type: actionTypes.auth.login.LOADING,
  };
}

export function loginSuccess(history, payload) {
  history.push("/");
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
  history.push("/logout");

  Cookies.expire(cookieNames.auth.token);
  Cookies.expire(cookieNames.auth.userId);
  Cookies.expire(cookieNames.auth.userEmail);
  Cookies.expire(cookieNames.auth.userProfileName);
  Cookies.expire(cookieNames.auth.username);

  return {
    type: actionTypes.auth.logout.REQUEST,
  };
}

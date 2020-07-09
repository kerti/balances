import actionTypes from "./actionTypes";
import auth from "../sources/auth/auth";
import cookieNames from "../cookies";

export function requestLogin(setCookies, username, password, history) {
  return (dispatch) => {
    dispatch(loginLoading());
    auth
      .login(username, password)
      .then((payload) => {
        const loginResponse = payload.data.data;

        setCookies(cookieNames.auth.token, loginResponse.token, {
          expires: new Date(loginResponse.expiration),
        });

        setCookies(cookieNames.auth.userId, loginResponse.user.id, {
          expires: new Date(payload.data.data.expiration),
        });

        setCookies(cookieNames.auth.userEmail, loginResponse.user.email, {
          expires: new Date(payload.data.data.expiration),
        });

        setCookies(cookieNames.auth.userProfileName, loginResponse.user.name, {
          expires: new Date(payload.data.data.expiration),
        });

        setCookies(cookieNames.auth.username, loginResponse.user.username, {
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

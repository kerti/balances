import { types } from "../../actions";

const auth = (state = {}, action) => {
  switch (action.type) {
    case types.auth.login.LOADING:
      return {
        loading: true,
      };
    case types.auth.login.SUCCESS:
    case types.auth.login.LOADCOOKIES:
      const data = action.payload.data;
      return {
        loading: false,
        token: data.token,
      };
    case types.auth.login.FAILURE:
      return {
        loading: false,
        authError: action.error,
      };
    case types.auth.logout.REQUEST:
      return {};
    default:
      return state;
  }
};

export default auth;

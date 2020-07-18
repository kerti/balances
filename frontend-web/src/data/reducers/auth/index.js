import { actionTypes } from "../../actions";

const auth = (state = {}, action) => {
  switch (action.type) {
    case actionTypes.auth.login.LOADING:
      return {
        loading: true,
      };
    case actionTypes.auth.login.SUCCESS:
    case actionTypes.auth.login.LOADCOOKIES:
      const data = action.payload.data;
      return {
        loading: false,
        token: data.token,
      };
    case actionTypes.auth.login.FAILURE:
      return {
        loading: false,
        authError: action.error,
      };
    case actionTypes.auth.logout.REQUEST:
      return {};
    default:
      return state;
  }
};

export default auth;

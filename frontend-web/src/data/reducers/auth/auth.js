import actionTypes from "../../actions/actionTypes";

const auth = (state = {}, action) => {
  switch (action.type) {
    case actionTypes.auth.login.LOADING:
      return {
        loading: true,
      };
    case actionTypes.auth.login.SUCCESS:
      return {
        loading: false,
        authData: action.data,
      };
    case actionTypes.auth.login.FAILURE:
      return {
        loading: false,
        authError: action.error,
      };
    default:
      return state;
  }
};

export default auth;

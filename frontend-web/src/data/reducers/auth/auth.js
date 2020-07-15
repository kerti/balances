import actionTypes from "../../actions/actionTypes";

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
        expiration: data.expiration,
        loading: false,
        token: data.token,
        user: {
          id: data.user.id,
          email: data.user.email,
          name: data.user.name,
          username: data.user.username,
        },
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

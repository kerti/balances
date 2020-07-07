const actionTypes = {
  auth: {
    login: {
      LOADING: "AUTH_LOGIN_LOADING",
      FAILURE: "AUTH_LOGIN_FAILURE",
      SUCCESS: "AUTH_LOGIN_SUCCESS",
    },
    logout: "AUTH_LOGOUT",
  },
  ui: {
    sidebarShow: {
      SET: "UI_SIDEBAR_SHOW_SET",
    },
  },
};

export default actionTypes;

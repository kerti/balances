const types = {
  auth: {
    login: {
      LOADING: "AUTH_LOGIN_LOADING",
      FAILURE: "AUTH_LOGIN_FAILURE",
      SUCCESS: "AUTH_LOGIN_SUCCESS",
      LOADCOOKIES: "AUTH_LOGIN_LOADCOOKIES",
    },
    logout: {
      REQUEST: "AUTH_LOGOUT_REQUEST",
    },
  },
  ui: {
    lang: {
      SET: "UI_LANG_SET",
    },
    sidebarShow: {
      SET: "UI_SIDEBAR_SHOW_SET",
    },
  },
};

export default types;

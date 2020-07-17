const actionTypes = {
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
  pages: {
    assets: {
      banks: {
        list: {
          LOADING: "PAGES_ASSETS_BANKS_LIST_LOADING",
          FAILURE: "PAGES_ASSETS_BANKS_LIST_FAILURE",
          SUCCESS: "PAGES_ASSETS_BANKS_LIST_SUCCESS",
        },
        detail: {
          LOADING: "PAGES_ASSETS_BANKS_DETAIL_LOADING",
          FAILURE: "PAGES_ASSETS_BANKS_DETAIL_FAILURE",
          SUCCESS: "PAGES_ASSETS_BANKS_DETAIL_SUCCESS",
        },
      },
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

export default actionTypes;

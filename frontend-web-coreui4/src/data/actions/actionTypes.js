const actionTypes = {
  auth: {
    login: {
      LOADING: 'AUTH_LOGIN_LOADING',
      FAILURE: 'AUTH_LOGIN_FAILURE',
      SUCCESS: 'AUTH_LOGIN_SUCCESS',
    },
    logout: {
      REQUEST: 'AUTH_LOGOUT_REQUEST',
    },
  },
  ui: {
    lang: {
      SET: 'UI_LANG_SET',
    },
    sidebarUnfoldable: {
      SET: 'UI_SIDEBAR_UNFOLDABLE_SET',
    },
    sidebarShow: {
      SET: 'UI_SIDEBAR_SHOW_SET',
    },
  },
}

export default actionTypes
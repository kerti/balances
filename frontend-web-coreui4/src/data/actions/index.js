const actionTypes = {
  auth: {
    login: {
      LOADING: 'AUTH_LOGIN_LOADING',
      FAILURE: 'AUTH_LOGIN_FAILURE',
      SUCCESS: 'AUTH_LOGIN_SUCCESS',
      LOADCOOKIES: 'AUTH_LOGIN_LOADCOOKIES',
    },
    logout: {
      REQUEST: 'AUTH_LOGOUT_REQUEST',
    },
  },
  entities: {
    user: {
      REQUEST: 'E_USER_REQUEST',
      SUCCESS: 'E_USER_SUCCESS',
      FAILURE: 'E_USER_FAILURE',
      page: {
        REQUEST: 'E_USER_PAGE_REQUEST',
        SUCCESS: 'E_USER_PAGE_SUCCESS',
        FAILURE: 'E_USER_PAGE_FAILURE',
      },
    },
  },
  ui: {
    errorMessage: {
      RESET: 'UI_ERROR_MESSAGE_RESET',
    },
    lang: {
      SET: 'UI_LANG_SET',
    },
    sidebarShow: {
      SET: 'UI_SIDEBAR_SHOW_SET',
    },
    sidebarUnfoldable: {
      SET: 'UI_SIDEBAR_UNFOLDABLE_SET',
    },
  },
}

export { actionTypes }

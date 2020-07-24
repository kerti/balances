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
  entities: {
    bankAccount: {
      REQUEST: "E_BANK_ACCOUNT_REQUEST",
      SUCCESS: "E_BANK_ACCOUNT_SUCCESS",
      FAILURE: "E_BANK_ACCOUNT_FAILURE",
      update: {
        REQUEST: "E_BANK_ACCOUNT_UPDATE_REQUEST",
        SUCCESS: "E_BANK_ACCOUNT_UPDATE_SUCCESS",
        FAILURE: "E_BANK_ACCOUNT_UPDATE_FAILURE",
      },
      page: {
        REQUEST: "E_BANK_ACCOUNT_PAGE_REQUEST",
        SUCCESS: "E_BANK_ACCOUNT_PAGE_SUCCESS",
        FAILURE: "E_BANK_ACCOUNT_PAGE_FAILURE",
      },
    },
    bankAccountBalance: {
      REQUEST: "E_BANK_ACCOUNT_BALANCE_REQUEST",
      SUCCESS: "E_BANK_ACCOUNT_BALANCE_SUCCESS",
      FAILURE: "E_BANK_ACCOUNT_BALANCE_FAILURE",
      update: {
        REQUEST: "E_BANK_ACCOUNT_BALANCE_UPDATE_REQUEST",
        SUCCESS: "E_BANK_ACCOUNT_BALANCE_UPDATE_SUCCESS",
        FAILURE: "E_BANK_ACCOUNT_BALANCE_UPDATE_FAILURE",
      },
      page: {
        REQUEST: "E_BANK_ACCOUNT_BALANCE_PAGE_REQUEST",
        SUCCESS: "E_BANK_ACCOUNT_BALANCE_PAGE_SUCCESS",
        FAILURE: "E_BANK_ACCOUNT_BALANCE_PAGE_FAILURE",
      },
    },
    user: {
      REQUEST: "E_USER_REQUEST",
      SUCCESS: "E_USER_SUCCESS",
      FAILURE: "E_USER_FAILURE",
      page: {
        REQUEST: "E_USER_PAGE_REQUEST",
        SUCCESS: "E_USER_PAGE_SUCCESS",
        FAILURE: "E_USER_PAGE_FAILURE",
      },
    },
  },
  ui: {
    errorMessage: {
      RESET: "UI_ERROR_MESSAGE_RESET",
    },
    lang: {
      SET: "UI_LANG_SET",
    },
    sidebarShow: {
      SET: "UI_SIDEBAR_SHOW_SET",
    },
  },
};

export { actionTypes };

import { actionTypes } from "../../actions";

import { CALL_API } from "../../middleware/api";
import { Schemas } from "../../schemas";

// Fetches a single Bank Acccount from Balances API.
// Relies on the custom API middlewre defined in ../middleware/api.js.
const fetchBankAccount = (id) => ({
  [CALL_API]: {
    types: [
      actionTypes.entities.bankAccount.REQUEST,
      actionTypes.entities.bankAccount.SUCCESS,
      actionTypes.entities.bankAccount.FAILURE,
    ],
    endpoint: `bankAccounts/${id}`,
    schema: Schemas.BANK_ACCOUNT,
  },
});

// Fetches a single Bank Account from Balances API unless it is cached.
// Relies on Redux Thunk middleware.
export const loadBankAccount = (id, requiredFields = []) => (
  dispatch,
  getState
) => {
  const bankAccount = getState().entities.bankAccounts[id];
  if (
    bankAccount &&
    requiredFields.every((key) => bankAccount.hasOwnProperty(key))
  ) {
    return null;
  }

  return dispatch(fetchBankAccount(id));
};

// Fetches a page of Bank Accounts for a particular keyword.
// Relies on the custom API middleware defined in ../middleware/api.js.
const fetchBankAccountPage = (
  keyword,
  page,
  pageSize = parseInt(process.env.REACT_APP_DEFAULT_PAGE_SIZE)
) => ({
  keyword,
  page,
  [CALL_API]: {
    types: [
      actionTypes.entities.bankAccountPage.REQUEST,
      actionTypes.entities.bankAccountPage.SUCCESS,
      actionTypes.entities.bankAccountPage.FAILURE,
    ],
    endpoint: `bankAccounts/search`,
    schema: Schemas.BANK_ACCOUNT_ARRAY,
    method: "POST",
    body: {
      keyword: keyword,
      page: page,
      pageSize: pageSize,
    },
  },
});

// Fetches a page of Bank Accounts for a particular keyword.
// Bails out if page is cached and user didn't specifically request next page.
// Relies on Redux Thunk middleware.
export const loadBankAccountPage = (
  keyword,
  page = 1,
  pageSize = parseInt(process.env.REACT_APP_DEFAULT_PAGE_SIZE)
) => (dispatch, getState) => {
  const { currentPage = 1, pageCount = 0 } =
    getState().pagination.bankAccountsByKeyword[keyword] || {};

  if (pageCount > 0 && page === currentPage) {
    return null;
  }

  return dispatch(fetchBankAccountPage(keyword, page, pageSize));
};

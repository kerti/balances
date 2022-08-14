import { actionTypes } from '../../actions'

import { CALL_API } from '../../middleware/api'
import { Schemas } from '../../schemas'

// Fetches a single Bank Acccount from Balances API.
// Relies on the custom API middleware defined in ../middleware/api.js.
const fetchBankAccount = (id, withBalances = true, balanceCount = 12) => ({
  [CALL_API]: {
    types: [
      actionTypes.entities.bankAccount.REQUEST,
      actionTypes.entities.bankAccount.SUCCESS,
      actionTypes.entities.bankAccount.FAILURE,
    ],
    endpoint: `bankAccounts/${id}?withBalances=${withBalances}&balanceCount=${balanceCount}`,
    schema: Schemas.BANK_ACCOUNT,
    method: 'GET',
  },
})

// Fetches a single Bank Account from Balances API unless it is cached.
// Relies on Redux Thunk middleware.
export const loadBankAccount =
  (
    id,
    withBalances = true,
    balanceCount = 12,
    requiredFields = [],
    ignoreCache = false
  ) =>
  (dispatch, getState) => {
    const bankAccount = getState().entities.bankAccounts[id]
    if (
      bankAccount &&
      bankAccount.balances &&
      bankAccount.balances.length > 0 &&
      requiredFields.every((key) => bankAccount.hasOwnProperty(key)) &&
      !ignoreCache
    ) {
      return null
    }

    return dispatch(fetchBankAccount(id, withBalances, balanceCount))
  }

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
      actionTypes.entities.bankAccount.page.REQUEST,
      actionTypes.entities.bankAccount.page.SUCCESS,
      actionTypes.entities.bankAccount.page.FAILURE,
    ],
    endpoint: `bankAccounts/search`,
    schema: Schemas.BANK_ACCOUNT_ARRAY,
    method: 'POST',
    body: {
      keyword: keyword,
      page: page,
      pageSize: pageSize,
    },
  },
})

// Fetches a page of Bank Accounts for a particular keyword.
// Bails out if page is cached and user didn't specifically request next page.
// Relies on Redux Thunk middleware.
export const loadBankAccountPage =
  (
    keyword,
    page = 1,
    pageSize = parseInt(process.env.REACT_APP_DEFAULT_PAGE_SIZE)
  ) =>
  (dispatch, getState) => {
    const { currentPage = 1, pageCount = 0 } =
      getState().pagination.bankAccountsByKeyword[keyword] || {}

    if (pageCount > 0 && page === currentPage) {
      return null
    }

    return dispatch(fetchBankAccountPage(keyword, page, pageSize))
  }

// Fetches a single Bank Account Balance from Balances API.
// Relies on the custom API middleware defined in ../middleware/api.js.
const fetchBankAccountBalance = (id) => ({
  [CALL_API]: {
    types: [
      actionTypes.entities.bankAccountBalance.REQUEST,
      actionTypes.entities.bankAccountBalance.SUCCESS,
      actionTypes.entities.bankAccountBalance.FAILURE,
    ],
    endpoint: `bankAccounts/balances/${id}`,
    schema: Schemas.BANK_ACCOUNT_BALANCE,
    method: 'GET',
  },
})

// Fetches a single Bank Account Balance from Balances API.
// Relies on Redux Thunk middleware.
export const loadBankAccountBalance =
  (id, requiredFields = []) =>
  (dispatch, getState) => {
    const bankAccountBalance = getState().entities.bankAccountBalances[id]
    if (
      bankAccountBalance &&
      requiredFields.every((key) => bankAccountBalance.hasOwnProperty(key))
    ) {
      return null
    }

    return dispatch(fetchBankAccountBalance(id))
  }

// Fetches a page of Bank Account Balances for a particular Bank Account.
// Relies on the custom API middleware defined in ../middleware/api.js.
const fetchBankAccountBalancePage = (
  bankAccountId,
  page,
  pageSize = parseInt(process.env.REACT_APP_DEFAULT_PAGE_SIZE)
) => ({
  bankAccountId,
  page,
  [CALL_API]: {
    types: [
      actionTypes.entities.bankAccountBalance.page.REQUEST,
      actionTypes.entities.bankAccountBalance.page.SUCCESS,
      actionTypes.entities.bankAccountBalance.page.FAILURE,
    ],
    endpoint: `bankAccounts/balances/search`,
    schema: Schemas.BANK_ACCOUNT_BALANCE_ARRAY,
    method: 'POST',
    body: {
      bankAccountId,
      page,
      pageSize,
    },
  },
})

// Fetches a page of Bank Account Balances for a particular Bank Account.
// Bails out if page is cached and user didn't specifically request next page
// or explicitly request to ignore the cache.
// Relies on Redux Thunk middleware.
export const loadBankAccountBalancePage =
  (
    bankAccountId,
    page,
    pageSize = parseInt(process.env.REACT_APP_DEFAULT_PAGE_SIZE),
    ignoreCache = false
  ) =>
  (dispatch, getState) => {
    const { currentPage = 1, pageCount = 0 } =
      getState().pagination.bankAccountBalancesByBankAccountId[bankAccountId] ||
      {}

    if (pageCount > 0 && page === currentPage && !ignoreCache) {
      return null
    }

    return dispatch(fetchBankAccountBalancePage(bankAccountId, page, pageSize))
  }

// Store a bank account.
export const updateBankAccount = (
  id,
  accountName,
  bankName,
  accountHolderName,
  accountNumber,
  status,
  options = {}
) => ({
  [CALL_API]: {
    types: [
      actionTypes.entities.bankAccount.update.REQUEST,
      actionTypes.entities.bankAccount.update.SUCCESS,
      actionTypes.entities.bankAccount.update.FAILURE,
    ],
    options: options,
    endpoint: `bankAccounts/${id}`,
    schema: Schemas.BANK_ACCOUNT,
    method: 'PATCH',
    body: {
      id,
      accountName,
      bankName,
      accountHolderName,
      accountNumber,
      status,
    },
  },
})

// Create new bank account balance.
export const createBankAccountBalance = (
  bankAccountId,
  date,
  balance,
  options = {}
) => ({
  [CALL_API]: {
    types: [
      actionTypes.entities.bankAccountBalance.create.REQUEST,
      actionTypes.entities.bankAccountBalance.create.SUCCESS,
      actionTypes.entities.bankAccountBalance.create.FAILURE,
    ],
    options: options,
    endpoint: `bankAccounts/balances`,
    schema: Schemas.BANK_ACCOUNT_BALANCE,
    method: 'POST',
    body: {
      bankAccountId,
      date,
      balance,
    },
  },
})

// Update existing bank balance.
export const updateBankAccountBalance = (
  bankAccountBalanceId,
  bankAccountId,
  date,
  balance,
  options = {}
) => ({
  [CALL_API]: {
    types: [
      actionTypes.entities.bankAccountBalance.update.REQUEST,
      actionTypes.entities.bankAccountBalance.update.SUCCESS,
      actionTypes.entities.bankAccountBalance.update.FAILURE,
    ],
    options: options,
    endpoints: 'bankAccounts/balances',
    schema: Schemas.BANK_ACCOUNT_BALANCE,
    method: 'PUT',
    body: {
      bankAccountBalanceId,
      bankAccountId,
      date,
      balance,
    },
  },
})

// Delete existing bank balance.
export const deleteBankAccountBalance = (
  bankAccountBalanceId,
  options = {}
) => ({
  [CALL_API]: {
    types: [
      actionTypes.entities.bankAccountBalance.delete.REQUEST,
      actionTypes.entities.bankAccountBalance.delete.SUCCESS,
      actionTypes.entities.bankAccountBalance.delete.FAILURE,
    ],
    options: options,
    endpoints: 'bankAccounts/balances',
    schema: Schemas.BANK_ACCOUNT_BALANCE,
    method: 'DELETE',
    body: {
      bankAccountBalanceId,
    },
  },
})

export const showBalanceModal = () => ({
  type: actionTypes.ui.modals.assets.bankAccounts.balances.SHOW,
})

export const hideBalanceModal = () => ({
  type: actionTypes.ui.modals.assets.bankAccounts.balances.HIDE,
})

import { combineReducers } from 'redux'
import auth from './auth'
import paginate from './paginate'
import apiState from './apiState'
import ui from './ui'
import { actionTypes } from '../actions'
import merge from 'lodash/merge'

// Updates an entity cache in response to any action with response.entities.
const entities = (
  state = {
    users: {},
    bankAccounts: {},
    bankAccountBalances: {},
  },
  action
) => {
  if (action.response && action.response.entities) {
    return merge({}, state, action.response.entities)
  }

  return state
}

// Updates error message to notify about the failed fetches.
const errorMessage = (state = null, action) => {
  const { type, error } = action

  if (type === actionTypes.ui.errorMessage.RESET) {
    return null
  } else if (error) {
    return error
  }

  return state
}

// Updates the pagination data for different actions.
const pagination = combineReducers({
  bankAccountsByKeyword: paginate({
    mapActionToKey: (action) => action.keyword,
    types: [
      actionTypes.entities.bankAccount.page.REQUEST,
      actionTypes.entities.bankAccount.page.SUCCESS,
      actionTypes.entities.bankAccount.page.FAILURE,
    ],
  }),
  bankAccountBalancesByBankAccountId: paginate({
    mapActionToKey: (action) => action.bankAccountId,
    types: [
      actionTypes.entities.bankAccountBalance.page.REQUEST,
      actionTypes.entities.bankAccountBalance.page.SUCCESS,
      actionTypes.entities.bankAccountBalance.page.FAILURE,
    ],
  }),
  usersByFilter: paginate({
    mapActionToKey: (action) => {
      const { ids, keyword, page, pageSize } = action
      const filter = { ids, keyword, page, pageSize }
      return JSON.stringify(filter)
    },
    types: [
      actionTypes.entities.user.page.REQUEST,
      actionTypes.entities.user.page.SUCCESS,
      actionTypes.entities.user.page.FAILURE,
    ],
  }),
})

// Updates the apiState data for different actions.
const api = combineReducers({
  // bank accounts
  updateBankAccount: apiState({
    types: [
      actionTypes.entities.bankAccount.update.REQUEST,
      actionTypes.entities.bankAccount.update.SUCCESS,
      actionTypes.entities.bankAccount.update.FAILURE,
    ],
  }),
  createBankAccountBalance: apiState({
    types: [
      actionTypes.entities.bankAccountBalance.create.REQUEST,
      actionTypes.entities.bankAccountBalance.create.SUCCESS,
      actionTypes.entities.bankAccountBalance.create.FAILURE,
    ],
  }),
})

const rootReducer = combineReducers({
  auth,
  entities,
  errorMessage,
  pagination,
  api,
  ui,
})

export default rootReducer

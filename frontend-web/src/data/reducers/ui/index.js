import { combineReducers } from 'redux'
import { actionTypes } from '../../actions'
import modalState from './modalState'

export function lang(state = 'en', action) {
  switch (action.type) {
    case actionTypes.ui.lang.SET:
      return action.data
    default:
      return state
  }
}

// Updates the modalState for different actions.
const modals = combineReducers({
  bankAccountsBalance: modalState({
    types: [
      actionTypes.ui.modals.assets.bankAccounts.balances.SHOW,
      actionTypes.ui.modals.assets.bankAccounts.balances.HIDE,
    ],
  }),
})

export function sidebarShow(state = 'responsive', action) {
  switch (action.type) {
    case actionTypes.ui.sidebarShow.SET:
      return action.data
    default:
      return state
  }
}

export default combineReducers({
  lang,
  modals,
  sidebarShow,
})

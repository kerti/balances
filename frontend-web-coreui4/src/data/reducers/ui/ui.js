import { combineReducers } from 'redux'

import actionTypes from '../../actions/actionTypes'

export function lang(state = 'en', action) {
  switch (action.type) {
    case actionTypes.ui.lang.SET:
      return action.data
    default:
      return state
  }
}

export function sidebarUnfoldable(state = false, action) {
  switch (action.type) {
    case actionTypes.ui.sidebarUnfoldable.SET:
      return action.data
    default:
      return state
  }
}

export function sidebarShow(state = true, action) {
  switch (action.type) {
    case actionTypes.ui.sidebarShow.SET:
      return action.data
    default:
      return state
  }
}

export default combineReducers({
  lang,
  sidebarUnfoldable,
  sidebarShow,
})

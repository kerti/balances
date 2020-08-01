import { combineReducers } from 'redux'

import { types } from '../../actions'

export function lang(state = 'en', action) {
  switch (action.type) {
    case types.ui.lang.SET:
      return action.data
    default:
      return state
  }
}

export function sidebarShow(state = 'responsive', action) {
  switch (action.type) {
    case types.ui.sidebarShow.SET:
      return action.data
    default:
      return state
  }
}

export default combineReducers({
  lang,
  sidebarShow,
})

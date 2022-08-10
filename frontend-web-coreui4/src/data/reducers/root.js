import { combineReducers } from 'redux'

import auth from './auth/auth'
import ui from './ui/ui'

export default combineReducers({
  auth: auth,
  ui: ui,
})

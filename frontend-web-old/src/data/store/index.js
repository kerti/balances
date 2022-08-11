import { createStore, applyMiddleware, compose } from 'redux'
import thunkMiddleware from 'redux-thunk'
import rootReducer from '../reducers'
import api from '../middleware/api'
import DevTools from '../../containers/DevTools'

const middlewares = [thunkMiddleware, api]

if (
  process.env.NODE_ENV === 'development' &&
  process.env.REACT_APP_USE_REDUX_LOGGER === 'true'
) {
  const { createLogger } = require('redux-logger')
  const logger = createLogger()
  middlewares.push(logger)
}

const initialState = {
  ui: {
    lang: process.env.REACT_APP_DEFAULT_LANG,
    sidebarShow: 'responsive',
  },
}

const initStore = () => {
  const store =
    process.env.NODE_ENV === 'development'
      ? createStore(
          rootReducer,
          initialState,
          compose(applyMiddleware(...middlewares), DevTools.instrument())
        )
      : createStore(rootReducer, initialState, applyMiddleware(...middlewares))

  if (module.hot) {
    module.hot.accept('../reducers', () => {
      const nextRootReducer = require('../reducers').default
      store.replaceReducer(nextRootReducer)
    })
  }

  return store
}

export default initStore

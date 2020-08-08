import 'react-app-polyfill/ie11' // For IE 11 support
import 'react-app-polyfill/stable'
import './polyfill'
import React from 'react'
import ReactDOM from 'react-dom'
import App from './App'
import * as serviceWorker from './serviceWorker'

import { icons } from './assets/icons'

import { Provider } from 'react-redux'
import initStore from './data/store'

import { I18nextProvider } from 'react-i18next'
import i18next from 'i18next'

import { initTranslations } from './translations'
import { setLangFromCookie } from './data/actions/ui'

import { loadAuthCookies } from './data/actions/system/auth'
import DevTools from './containers/DevTools'

initTranslations()

React.icons = icons

const store = initStore()

store.dispatch(setLangFromCookie())
store.dispatch(loadAuthCookies())

// TODO: don't include devtools in production
ReactDOM.render(
  <Provider store={store}>
    <I18nextProvider i18n={i18next}>
      <App />
      <DevTools />
    </I18nextProvider>
  </Provider>,
  document.getElementById('root')
)

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: http://bit.ly/CRA-PWA
serviceWorker.unregister()

import 'react-app-polyfill/stable'
import 'core-js'
import React from 'react'
import ReactDOM from 'react-dom'
import App from './App'
import reportWebVitals from './reportWebVitals'
import { Provider } from 'react-redux'
import initStore from './data/store'
import { I18nextProvider } from 'react-i18next'
import i18next from 'i18next'
import { initTranslations } from './translations'
import { setLangFromCookie } from './data/actions/ui'

import { loadAuthCookies } from './data/actions/system/auth'
import DevTools from './components/DevTools'

initTranslations()

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
  document.getElementById('root'),
)

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals()

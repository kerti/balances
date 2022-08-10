import 'react-app-polyfill/stable'
import 'core-js'
import React from 'react'
import ReactDOM from 'react-dom'
import App from './App'
import reportWebVitals from './reportWebVitals'
import { Provider } from 'react-redux'
import initStore from './data/store'
import { CookiesProvider } from 'react-cookie'
import { I18nextProvider } from 'react-i18next'
import i18next from 'i18next'
import app_en from './translations/en/app.json'
import app_id from './translations/id/app.json'

const store = initStore()

i18next.init({
  interpolation: { escapeValue: false },
  lng: process.env.REACT_APP_DEFAULT_LANG,
  resources: {
    en: {
      app: app_en,
    },
    id: {
      app: app_id,
    },
  },
})

ReactDOM.render(
  <CookiesProvider>
    <Provider store={store}>
      <I18nextProvider i18n={i18next}>
        <App />
      </I18nextProvider>
    </Provider>
  </CookiesProvider>,
  document.getElementById('root'),
)

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals()

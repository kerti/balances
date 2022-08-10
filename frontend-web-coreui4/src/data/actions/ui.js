import types from './types'
import i18n from 'i18next'
import Cookies from 'cookies-js'
import cookieNames from '../cookies'

export function setLangFromCookie() {
  const lang = Cookies.get(cookieNames.ui.lang) || process.env.REACT_APP_DEFAULT_LANG
  i18n.changeLanguage(lang)
  return {
    type: types.ui.lang.SET,
    data: lang,
  }
}

export function setLang(data) {
  Cookies.set(cookieNames.ui.lang, data)
  i18n.changeLanguage(data)
  return {
    type: types.ui.lang.SET,
    data,
  }
}

export function setSidebarUnfoldable(data) {
  return {
    type: types.ui.sidebarUnfoldable.SET,
    data,
  }
}

export function setSidebarShow(data) {
  return {
    type: types.ui.sidebarShow.SET,
    data,
  }
}

import actionTypes from './actionTypes'

export function setLang(data) {
  return {
    type: actionTypes.ui.lang.SET,
    data,
  }
}

export function setSidebarUnfoldable(data) {
  return {
    type: actionTypes.ui.sidebarUnfoldable.SET,
    data,
  }
}

export function setSidebarShow(data) {
  return {
    type: actionTypes.ui.sidebarShow.SET,
    data,
  }
}

import actionTypes from "./actionTypes";

export function setSidebarShow(data) {
  return {
    type: actionTypes.ui.sidebarShow.SET,
    data,
  };
}

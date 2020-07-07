import actionTypes from "../../actions/actionTypes";

export function sidebarShow(state = "responsive", action) {
  switch (action.type) {
    case actionTypes.ui.sidebarShow.SET:
      return action.data;
    default:
      return state;
  }
}

import banks from "./banks/banks";
import { combineReducers } from "redux";

export default combineReducers({
  banks: banks,
});

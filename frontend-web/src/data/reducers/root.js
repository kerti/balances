import { combineReducers } from "redux";

import auth from "./auth/auth";
import ui from "./ui/ui";
import pages from "./pages/pages";

export default combineReducers({
  auth: auth,
  pages: pages,
  ui: ui,
});

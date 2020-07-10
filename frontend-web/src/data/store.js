import { createStore, applyMiddleware } from "redux";
import thunkMiddleware from "redux-thunk";
import rootReducer from "./reducers/root";

const middlewares = [thunkMiddleware];

if (process.env.NODE_ENV === "development") {
  const { createLogger } = require("redux-logger");
  const logger = createLogger();
  middlewares.push(logger);
}

const initialState = {
  ui: {
    lang: process.env.REACT_APP_DEFAULT_LANG,
    sidebarShow: "responsive",
  },
};

const initStore = () => {
  const store = createStore(
    rootReducer,
    initialState,
    applyMiddleware(...middlewares)
  );

  if (module.hot) {
    module.hot.accept("./reducers/root", () => {
      const nextRootReducer = require("./reducers/root").default;
      store.replaceReducer(nextRootReducer);
    });
  }

  return store;
};

export default initStore;

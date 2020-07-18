import { combineReducers } from "redux";
import merge from "lodash/merge";
import auth from "./auth";
import paginate from "./paginate";
import ui from "./ui";
import { actionTypes } from "../actions";

// Updates an entity cache in response to any action with response.entities.
const entities = (
  state = { users: {}, bankAccounts: {}, bankAccountBalances: {} },
  action
) => {
  if (action.response && action.response.entities) {
    return merge({}, state, action.response.entities);
  }

  return state;
};

// Updates error message to notify about the failed fetches.
const errorMessage = (state = null, action) => {
  const { type, error } = action;

  if (type === actionTypes.ui.errorMessage.RESET) {
    return null;
  } else if (error) {
    return error;
  }

  return state;
};

// Updates the pagination data for different actions.
const pagination = combineReducers({
  bankAccountsByKeyword: paginate({
    mapActionToKey: (action) => action.keyword,
    types: [
      actionTypes.entities.bankAccountPage.REQUEST,
      actionTypes.entities.bankAccountPage.SUCCESS,
      actionTypes.entities.bankAccountPage.FAILURE,
    ],
  }),
});

const rootReducer = combineReducers({
  auth,
  entities,
  errorMessage,
  pagination,
  ui,
});

export default rootReducer;

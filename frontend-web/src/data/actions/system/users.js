import { actionTypes } from "../../actions";

import { CALL_API, Schemas } from "../middleware/api";

// Fetches a single user from Balances API.
// Relies on the custom API middleware defined in ../middleware/api.js.
const fetchUser = (id) => ({
  [CALL_API]: {
    types: [
      actionTypes.entities.user.REQUEST,
      actionTypes.entities.user.SUCCESS,
      actionTypes.entities.user.FAILURE,
    ],
    endpoint: `users/${id}`,
    schema: Schemas.USER,
  },
});

// Fetches a single user from Balances API unless it is cached.
// Relies on Redux Thunk middleware.
export const loadUser = (id, requiredFields = []) => (dispatch, getState) => {
  const user = getState().entities.users[id];
  if (user && requiredFields.every((key) => user.hasOwnProperty(key))) {
    return null;
  }

  return dispatch(fetchUser(id));
};

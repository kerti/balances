import { actionTypes } from "../../actions";

import { CALL_API } from "../../middleware/api";
import { Schemas } from "../../schemas";

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

// Fetches a page of Users based on a filter.
// Relies on the custom API middleware defined in ../middleware/api.js.
const fetchUserPage = (
  ids,
  keyword,
  page,
  pageSize = parseInt(process.env.REACT_APP_DEFAULT_PAGE_SIZE)
) => ({
  ids,
  keyword,
  page,
  [CALL_API]: {
    types: [
      actionTypes.entities.user.page.REQUEST,
      actionTypes.entities.user.page.SUCCESS,
      actionTypes.entities.user.page.FAILURE,
    ],
    endpoint: `users/search`,
    schema: Schemas.USER_ARRAY,
    method: "POST",
    body: {
      ids: ids,
      keyword: keyword,
      page: page,
      pageSize: pageSize,
    },
  },
});

// Fetches a page of Users based on a filter.
// Bails out if page is cached and user didn't specifically request based on different filter.
// Relies on Redux Thunk middleware.
export const loadUserPage = (
  ids,
  keyword,
  page = 1,
  pageSize = parseInt(process.env.REACT_APP_DEFAULT_PAGE_SIZE)
) => (dispatch, getState) => {
  const filter = { ids, keyword, page, pageSize };
  const filterString = JSON.stringify(filter);
  const { currentPage = 1, pageCount = 0 } =
    getState().pagination.usersByFilter[filterString] || {};

  if (pageCount > 0 && page === currentPage) {
    return null;
  }

  return dispatch(fetchUserPage(ids, keyword, page, pageSize));
};

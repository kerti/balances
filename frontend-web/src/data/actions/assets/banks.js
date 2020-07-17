import actionTypes from "../actionTypes";
import banks from "../../sources/assets/banks";

export function requestBankList(
  keyword,
  page = 1,
  pageSize = process.env.REACT_APP_DEFAULT_PAGE_SIZE,
  includeDeleted = false
) {
  return (dispatch) => {
    dispatch(bankSearchLoading());
    banks
      .getBankList(keyword, page, pageSize, includeDeleted)
      .then((payload) => {
        dispatch(bankSearchSuccess(payload.data.data));
      })
      .catch((error) => {
        dispatch(bankSearchError(error.response.data));
      });
  };
}

export function bankSearchLoading() {
  return {
    type: actionTypes.pages.assets.banks.list.LOADING,
  };
}

export function bankSearchSuccess(payload) {
  return {
    type: actionTypes.pages.assets.banks.list.SUCCESS,
    payload: payload,
  };
}

export function bankSearchError(error) {
  return {
    type: actionTypes.pages.assets.banks.list.FAILURE,
    error: error,
  };
}

export function requestBankDetailByID(
  id,
  withBalances = true,
  balanceCount = process.env.REACT_APP_DEFAULT_PAGE_SIZE
) {
  return (dispatch) => {
    dispatch(bankDetailLoading());
    banks
      .getBankByID(id, withBalances, balanceCount)
      .then((payload) => {
        dispatch(bankDetailSuccess(payload.data.data));
      })
      .catch((error) => {
        dispatch(bankDetailError(error.response.data));
      });
  };
}

export function bankDetailLoading() {
  return {
    type: actionTypes.pages.assets.banks.detail.LOADING,
  };
}

export function bankDetailSuccess(payload) {
  return {
    type: actionTypes.pages.assets.banks.detail.SUCCESS,
    payload: payload,
  };
}

export function bankDetailError(error) {
  return {
    type: actionTypes.pages.assets.banks.list.FAILURE,
    error: error,
  };
}

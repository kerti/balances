import actionTypes from "../../../../actions/actionTypes";

const banks = (state = {}, action) => {
  switch (action.type) {
    case actionTypes.pages.assets.banks.list.LOADING:
    case actionTypes.pages.assets.banks.detail.LOADING:
      return {
        loading: true,
      };
    case actionTypes.pages.assets.banks.list.FAILURE:
    case actionTypes.pages.assets.banks.detail.FAILURE:
      return {
        loading: false,
        error: action.error,
      };
    case actionTypes.pages.assets.banks.list.SUCCESS:
      return {
        loading: false,
        items: action.payload.items,
        pageInfo: action.payload.pageInfo,
      };
    case actionTypes.pages.assets.banks.detail.SUCCESS:
      return {
        loading: false,
        bank: action.payload,
      };
    default:
      return state;
  }
};

export default banks;

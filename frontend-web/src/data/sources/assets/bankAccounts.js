import axios from "axios";
import sources from "../sources";

export default {
  getBankAccountList: (
    keyword,
    page = 1,
    pageSize = process.env.REACT_APP_DEFAULT_PAGE_SIZE,
    includeDeleted = false
  ) => {
    return axios.post(
      `${sources.baseURL}/bankAccounts/search`,
      {
        keyword: keyword,
        includeDeleted: includeDeleted,
        page: parseInt(page),
        pageSize: parseInt(pageSize),
      },
      { withCredentials: true }
    );
  },

  getBankAccountByID: (id, withBalances = true, balanceCount = 12) => {
    return withBalances
      ? axios.get(
          `${sources.baseURL}/bankAccounts/${id}?withBalances&balanceCount=${balanceCount}`,
          {
            withCredentials: true,
          }
        )
      : axios.get(`${sources.baseURL}/bankAccounts/${id}`, {
          withCredentials: true,
        });
  },
};

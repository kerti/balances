import http from "axios";
import sources from "../sources";

export default {
  login: (username, password) => {
    const basic = btoa(`${username}:${password}`);
    return http.post(`${sources.baseURL}/auth/login`, null, {
      headers: {
        Authorization: `Basic ${basic}`,
      },
    });
  },
};

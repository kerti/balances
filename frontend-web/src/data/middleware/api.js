import { normalize } from "normalizr";
import { camelizeKeys } from "humps";
import axios from "axios";
import sources from "../sources";

// Fetches an API response and normalizes the result JSON according to schema.
// This makes every API response have the same shape, regardless of how nested it was.
const callApi = (endpoint, schema, options = {}) => {
  const fullUrl =
    endpoint.indexOf(sources.baseURL) === -1
      ? sources.baseURL + endpoint
      : endpoint;

  return axios({
    ...options,
    url: fullUrl,
    withCredentials: true,
  })
    .then((response) => {
      const camelizedJson = camelizeKeys(response.data.data.items);
      return Object.assign({}, normalize(camelizedJson, schema));
    })
    .catch((error) => {
      return Promise.reject(error);
    });
};

// Action key that carries API call info interpreted by this Redux middleware.
export const CALL_API = "Call API";

// A Redux middleware that interprets actions with CALL_API info specified.
// Performs the call and promises when such actions are dispatched.
export default (store) => (next) => (action) => {
  const callAPI = action[CALL_API];
  if (typeof callAPI === "undefined") {
    return next(action);
  }

  let { endpoint } = callAPI;
  const { schema, types, method, body } = callAPI;

  if (typeof endpoint === "function") {
    endpoint = endpoint(store.getState());
  }

  if (typeof endpoint !== "string") {
    throw new Error("Specify a string endpoint URL.");
  }

  if (!schema) {
    throw new Error("Specify one of the exported Schemas.");
  }

  if (!Array.isArray(types) || types.length !== 3) {
    throw new Error("Expected an array of three action types.");
  }

  if (!types.every((type) => typeof type === "string")) {
    throw new Error("Expected action types to be strings.");
  }

  const actionWith = (data) => {
    const finalAction = Object.assign({}, action, data);
    delete finalAction[CALL_API];
    return finalAction;
  };

  const [requestType, successType, failureType] = types;
  next(actionWith({ type: requestType }));

  const options = {
    method: method || "GET",
    data: body,
  };

  return callApi(endpoint, schema, options).then(
    (response) =>
      next(
        actionWith({
          response,
          type: successType,
        })
      ),
    (error) =>
      next(
        actionWith({
          type: failureType,
          error: error.message || "Something bad happened",
        })
      )
  );
};

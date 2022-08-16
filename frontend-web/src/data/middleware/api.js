import { normalize } from 'normalizr'
import { camelizeKeys } from 'humps'
import axios from 'axios'
import sources from '../sources'

const getPageInfo = (response) => {
  const data = response.data.data
  if (data.pageInfo) {
    return {
      pageCount: data.pageInfo.pageCount,
      totalCount: data.pageInfo.totalCount,
    }
  }
  return { pageCount: null, totalCount: null }
}

const getItems = (response) => {
  const data = response.data.data
  if (data.items && Array.isArray(data.items)) {
    return data.items
  }
  return data
}

// Fetches an API response and normalizes the result JSON according to schema.
// This makes every API response have the same shape, regardless of how nested it was.
const callApi = (endpoint, schema, options = {}) => {
  const fullUrl =
    endpoint.indexOf(sources.baseURL) === -1
      ? sources.baseURL + endpoint
      : endpoint

  return axios({
    ...options,
    url: fullUrl,
    withCredentials: true,
  })
    .then((response) => {
      const { pageCount, totalCount } = getPageInfo(response)
      const camelizedJson = camelizeKeys(getItems(response))
      const normalized = normalize(camelizedJson, schema)

      return Object.assign({}, normalized, {
        pageCount,
        totalCount,
      })
    })
    .catch((error) => {
      return Promise.reject(error)
    })
}

// Action key that carries API call info interpreted by this Redux middleware.
export const CALL_API = 'Call API'

// A Redux middleware that interprets actions with CALL_API info specified.
// Performs the call and promises when such actions are dispatched.
const api = (store) => (next) => (action) => {
  const callAPI = action[CALL_API]
  if (typeof callAPI === 'undefined') {
    return next(action)
  }

  let { endpoint } = callAPI
  const { schema, types, method, body, options } = callAPI
  let { nextAction } = options || {}

  if (typeof endpoint === 'function') {
    endpoint = endpoint(store.getState())
  }

  if (typeof endpoint !== 'string') {
    throw new Error('Specify a string endpoint URL.')
  }

  if (!schema) {
    throw new Error('Specify one of the exported Schemas.')
  }

  if (!Array.isArray(types) || types.length !== 3) {
    throw new Error('Expected an array of three action types.')
  }

  if (!types.every((type) => typeof type === 'string')) {
    throw new Error('Expected action types to be strings.')
  }

  const actionWith = (data) => {
    const finalAction = Object.assign({}, action, data)
    delete finalAction[CALL_API]
    return finalAction
  }

  if (nextAction !== undefined) {
    if (typeof nextAction === 'string') {
      nextAction = actionWith({ type: nextAction })
    }

    if (typeof nextAction === 'function') {
      nextAction = nextAction()
    }
  }

  const [requestType, successType, failureType] = types
  next(actionWith({ type: requestType }))

  const apiOptions = {
    method: method || 'GET',
    data: body,
  }

  return callApi(endpoint, schema, apiOptions).then(
    (response) => {
      next(
        actionWith({
          response,
          type: successType,
        })
      )
      if (nextAction !== undefined) {
        next(nextAction)
      }
    },
    (error) =>
      next(
        actionWith({
          type: failureType,
          error: error.message || 'Something bad happened',
        })
      )
  )
}

export default api

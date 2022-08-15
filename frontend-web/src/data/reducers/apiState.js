// Creates a reducer managing apiState, given the action types to handle.
const apiState = ({ types }) => {
  if (!Array.isArray(types) || types.length !== 3) {
    throw new Error('Expected types to be an array of three elements.')
  }
  if (!types.every((t) => typeof t === 'string')) {
    throw new Error('Expected types to be strings.')
  }

  const [requestType, successType, failureType] = types

  return (
    state = {
      isFetching: false,
    },
    action
  ) => {
    switch (action.type) {
      case requestType:
        return {
          ...state,
          isFetching: true,
        }
      case successType:
        return {
          ...state,
          isFetching: false,
        }
      case failureType:
        return {
          ...state,
          isFetching: false,
        }
      default:
        return state
    }
  }
}

export default apiState

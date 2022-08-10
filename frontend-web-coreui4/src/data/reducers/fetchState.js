// Creates a reducer managing fetchState, given the action types to handle,
// and a function telling how to extract the key from an action.
const fetchState = ({ types }) => {
  if (!Array.isArray(types) || types.length !== 3) {
    throw new Error('Expected types to be an array of three elements.')
  }
  if (!types.every((t) => typeof t === 'string')) {
    throw new Error('Expected types to be strings.')
  }

  const [requestType, successType, failureType] = types

  return (state = false, action) => {
    switch (action.type) {
      case requestType:
        return true
      case successType:
      case failureType:
        return false
      default:
        return state
    }
  }
}

export default fetchState

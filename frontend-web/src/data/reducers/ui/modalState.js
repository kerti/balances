// Creates a reducer managing UI modalState, given the actions to handle.
const modalState = ({ types }) => {
  if (!Array.isArray(types) || types.length !== 2) {
    throw new Error('Expected types to be an array of two elements.')
  }
  if (!types.every((t) => typeof t === 'string')) {
    throw new Error('Expected types to be strings.')
  }

  const [showType, hideType] = types

  return (
    state = {
      show: false,
    },
    action
  ) => {
    switch (action.type) {
      case showType:
        return {
          ...state,
          show: true,
        }
      case hideType:
        return {
          ...state,
          show: false,
        }
      default:
        return state
    }
  }
}

export default modalState

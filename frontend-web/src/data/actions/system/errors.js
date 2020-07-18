import { actionTypes } from "../../actions";

// Resets the currently visible error message.
export const resetErrorMessage = () => ({
  type: actionTypes.ui.errorMessage.RESET,
});

import { auth as actions } from '../actions';
import { auth as initialState } from '../state';

const {
  POST_SIGNIN_REQUEST,
  POST_SIGNIN_SUCCESS,
  POST_SIGNIN_FAILURE,
  POST_SIGNUP_REQUEST,
  POST_SIGNUP_SUCCESS,
  POST_SIGNUP_FAILURE,
  SIGN_OUT,
} = actions;

const reducer = (state = initialState, action) => {
  switch (action.type) {
    case POST_SIGNIN_REQUEST:
    case POST_SIGNUP_REQUEST:
      return { ...state, isFetching: true };
    case POST_SIGNIN_SUCCESS:
    case POST_SIGNUP_SUCCESS:
      return { ...state, isFetching: false, isLoggedIn: true, token: action.payload.token };
    case POST_SIGNIN_FAILURE:
    case POST_SIGNUP_FAILURE:
      return { ...state, isFetching: false, isLoggedIn: false };
    case SIGN_OUT:
      return { ...state, isLoggedIn: false, token: '' }
    default:
      return state;
  }
}

export default reducer;

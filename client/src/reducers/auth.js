import { auth } from '../actions';

const {
  POST_SIGNIN_REQUEST,
  POST_SIGNIN_SUCCESS,
  POST_SIGNIN_FAILURE,
  POST_SIGNUP_REQUEST,
  POST_SIGNUP_SUCCESS,
  POST_SIGNUP_FAILURE,
} = auth;

const initialState = {
  isFetching: false,
  isLoggedIn: false,
  token: '',
};

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
    default:
      return state;
  }
}

export default reducer;

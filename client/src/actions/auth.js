export const POST_SIGNUP_REQUEST = 'POST_SIGNUP_REQUEST';
export const POST_SIGNUP_SUCCESS = 'POST_SIGNUP_SUCCESS';
export const POST_SIGNUP_FAILURE = 'POST_SIGNUP_FAILURE';

export const postSignupRequest = () => ({
  type: POST_SIGNUP_REQUEST,
});

export const postSignupSuccess = payload => ({
  type: POST_SIGNUP_SUCCESS,
  payload,
});

export const postSignupFailure = payload => ({
  type: POST_SIGNUP_FAILURE,
  payload,
});

export const postSignup = (body) => async (dispatch) => {
  dispatch(postSignupRequest());
  try {
    const init = {
      method: 'POST',
      body: JSON.stringify(body),
    };
    const response = await fetch('/api/signup', init);
    const payload = await response.json();

    if (response.ok) {
      return dispatch(postSignupSuccess(payload));
    } else {
      throw dispatch(postSignupFailure(payload));
    }
  } catch (e) {
    throw e;
  }
};

import { applyMiddleware, createStore } from 'redux';
import { composeWithDevTools } from 'redux-devtools-extension';
import { createLogger } from 'redux-logger';
import thunk from 'redux-thunk';

import reducer from './reducers';
import state from './state';

const token = localStorage.getItem('token');
const isLoggedIn = Boolean(token);

const logger = createLogger();
const store = createStore(
  reducer,
  {
    ...state,
    auth: {
      ...state.auth,
      isLoggedIn,
      token,
    },
  },
  composeWithDevTools(
    applyMiddleware(thunk, logger),
  )
);

export default store;

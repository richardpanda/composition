import { applyMiddleware, compose, createStore } from 'redux';
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
  compose(
    applyMiddleware(thunk, logger),
    window.__REDUX_DEVTOOLS_EXTENSION__ && window.__REDUX_DEVTOOLS_EXTENSION__(),
  )
);

export default store;

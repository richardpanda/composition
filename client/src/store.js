import createHistory from 'history/createBrowserHistory';
import { routerMiddleware } from 'react-router-redux';
import { applyMiddleware, createStore } from 'redux';
import { composeWithDevTools } from 'redux-devtools-extension';
import { createLogger } from 'redux-logger';
import thunk from 'redux-thunk';

import reducers from './reducers';
import state from './state';

const history = createHistory();
const token = localStorage.getItem('token');
const isLoggedIn = Boolean(token);
const logger = createLogger();

const initialState = {
  ...state,
  auth: {
    ...state.auth,
    isLoggedIn,
    token,
  },
};
const middlewares = composeWithDevTools(
  applyMiddleware(thunk, logger, routerMiddleware(history)),
);

const store = createStore(
  reducers,
  initialState,
  middlewares,
);

export default store;

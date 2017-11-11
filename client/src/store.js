import createHistory from 'history/createBrowserHistory';
import jwtDecode from 'jwt-decode';
import { routerMiddleware } from 'react-router-redux';
import { applyMiddleware, createStore } from 'redux';
import { composeWithDevTools } from 'redux-devtools-extension';
import { createLogger } from 'redux-logger';
import thunk from 'redux-thunk';

import reducers from './reducers';
import state from './state';

const history = createHistory();
const token = localStorage.getItem('token') || '';
const isLoggedIn = Boolean(token);
const username = token ? jwtDecode(token).username : '';
const logger = createLogger();

const initialState = {
  ...state,
  auth: {
    ...state.auth,
    isLoggedIn,
    token,
    username,
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

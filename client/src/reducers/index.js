import { routerReducer } from 'react-router-redux';
import { combineReducers } from 'redux';

import auth from './auth';

export default combineReducers({
  auth,
  router: routerReducer,
});

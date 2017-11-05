import React, { Component } from 'react';
import {
  BrowserRouter as Router,
  Route,
  Switch,
} from 'react-router-dom';

import Navbar from './containers/Navbar';
import Signin from './containers/Signin';
import Signout from './containers/Signout';
import Signup from './containers/Signup';

class App extends Component {
  render() {
    return (
      <Router>
        <div>
          <header>
            <Navbar />
          </header>
          <main className="container">
            <Switch>
              <Route path="/signin" component={Signin} />
              <Route path="/signout" component={Signout} />
              <Route path="/signup" component={Signup} />
            </Switch>
          </main>
        </div>
      </Router>
    );
  }
}

export default App;

import React, { Component } from 'react';
import {
  Route,
  Switch,
} from 'react-router-dom';

import Home from './containers/Home';
import Navbar from './containers/Navbar';
import NewArticle from './containers/NewArticle';
import Signin from './containers/Signin';
import Signout from './containers/Signout';
import Signup from './containers/Signup';

class App extends Component {
  render() {
    return (
      <div>
        <header className="mb-4">
          <Navbar />
        </header>
        <main className="container">
          <Switch>
            <Route exact path="/" component={Home} />
            <Route path="/articles/new" component={NewArticle} />
            <Route path="/signin" component={Signin} />
            <Route path="/signout" component={Signout} />
            <Route path="/signup" component={Signup} />
          </Switch>
        </main>
      </div>
    );
  }
}

export default App;

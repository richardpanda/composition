import React, { Component } from 'react';
import { BrowserRouter as Router } from 'react-router-dom';

import Navbar from './components/Navbar';

class App extends Component {
  render() {
    return (
      <Router>
        <div>
          <header>
            <Navbar />
          </header>
        </div>
      </Router>
    );
  }
}

export default App;

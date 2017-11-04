import React, { Component } from 'react';
import { Link } from 'react-router-dom';

import './style.css';

class Navbar extends Component {
  render() {
    return (
      <nav className="navbar navbar-dark bg-dark">
        <a className="navbar-brand" href="/">Composition</a>
        <div>
          <Link to="/signin"><button className="auth-button btn btn-primary mr-2">Sign In</button></Link>
          <Link to="/signup"><button className="auth-button btn btn-light">Sign Up</button></Link>
        </div>
      </nav>
    );
  }
}

export default Navbar;

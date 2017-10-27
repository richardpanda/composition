import React, { Component } from 'react';

import './style.css';

class Navbar extends Component {
  render() {
    return (
      <nav className="Navbar-nav">
        <a className="Navbar-home" href="/">Composition</a>
        <ul className="Navbar-list">
          <li className="Navbar-item"><a className="Navbar-link" href="/signin">Sign In</a></li>
          <li className="Navbar-item"><a className="Navbar-link" href="/signup">Sign Up</a></li>
        </ul>
      </nav>
    );
  }
}

export default Navbar;

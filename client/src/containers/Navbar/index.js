import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';

import './style.css';

class Navbar extends Component {
  render() {
    const { isLoggedIn } = this.props;

    return (
      <nav className="navbar navbar-dark bg-dark">
        <a className="navbar-brand" href="/">Composition</a>
        {isLoggedIn ? (
          <Link to="/signout"><button className="auth-button btn btn-primary">Sign Out</button></Link>
        ) : (
          <div>
            <Link to="/signin"><button className="auth-button btn btn-primary mr-2">Sign In</button></Link>
            <Link to="/signup"><button className="auth-button btn btn-light">Sign Up</button></Link>
          </div>
        )}
      </nav>
    );
  }
}

const mapStateToProps = state => ({
  isLoggedIn: state.auth.isLoggedIn,
});

export default connect(mapStateToProps)(Navbar);
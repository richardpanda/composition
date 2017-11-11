import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';

class Navbar extends Component {
  render() {
    const { isLoggedIn, username } = this.props;

    return (
      <nav className="navbar navbar-dark navbar-expand bg-dark">
        <Link to="/"><div className="navbar-brand">Composition</div></Link>
        <div className="ml-auto">
          {isLoggedIn ? (
            <div className="navbar px-0 py-0">
              <ul className="navbar-nav">
                <li className="nav-item dropdown">
                  <a className="nav-link dropdown-toggle text-primary" href="#" data-toggle="dropdown">
                    {username}
                  </a>
                  <div className="dropdown-menu dropdown-menu-right">
                    <Link className="dropdown-item" to="/signout">Sign Out</Link>
                  </div>
                </li>
              </ul>
            </div>
          ) : (
            <div>
              <Link to="/signin"><button className="btn btn-link mr-2">Sign In</button></Link>
              <Link to="/signup"><button className="btn btn-light">Sign Up</button></Link>
            </div>
          )}
        </div>
      </nav>
    );
  }
}

const mapStateToProps = state => ({
  isLoggedIn: state.auth.isLoggedIn,
  username: state.auth.username,
});

export default connect(mapStateToProps)(Navbar);

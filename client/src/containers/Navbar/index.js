import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';

class Navbar extends Component {
  render() {
    const { isLoggedIn } = this.props;

    return (
      <nav className="navbar navbar-dark bg-dark">
        <Link to="/"><div className="navbar-brand">Composition</div></Link>
        {isLoggedIn ? (
          <Link to="/signout"><button className="btn btn-primary">Sign Out</button></Link>
        ) : (
          <div>
            <Link to="/signin"><button className="btn btn-primary mr-2">Sign In</button></Link>
            <Link to="/signup"><button className="btn btn-light">Sign Up</button></Link>
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

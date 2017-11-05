import React, { Component } from 'react';
import { connect } from 'react-redux';

import { auth as actions } from '../../actions';

const { signOut } = actions;

class Signout extends Component {
  componentDidMount() {
    const { history, signOut } = this.props;
    signOut();
    localStorage.removeItem('token');
    history.replace('/');
  }

  render() {
    return (
      <div>Signing out...</div>
    );
  }
}

const mapDispatchToProps = dispatch => ({
  signOut: () => dispatch(signOut()),
});

export default connect(null, mapDispatchToProps)(Signout);

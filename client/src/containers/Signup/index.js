import React, { Component } from 'react';
import { connect } from 'react-redux';

import { auth } from '../../actions';

const { postSignup } = auth;

class Signup extends Component {
  constructor(props) {
    super(props);

    this.state = {
      username: '',
      email: '',
      password: '',
      password_confirm: '',
      error: '',
    };
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(event) {
    const { name, value } = event.target;
    this.setState({
      [name]: value,
    });
  }

  async handleSubmit(event) {
    event.preventDefault();

    const { history, postSignup } = this.props;
    const { username, email, password, password_confirm } = this.state;
    const body = {
      username,
      email,
      password,
      password_confirm,
    };

    try {
        const action = await postSignup(body);
        localStorage.setItem('token', action.payload.token);
        history.push('/');
    } catch (e) {
        this.setState({ error: e.payload.message });
    }
  }

  render() {
    const { error } = this.state;

    return (
      <form className="w-75 mx-auto" onSubmit={this.handleSubmit}>
        <h3 className="text-center">Create an account</h3>
        {error && <div className="alert alert-danger">{error}</div>}
        <div className="form-group">
          <label>Username</label>
          <input type="text" className="form-control" name="username" onChange={this.handleChange} required />
        </div>
        <div className="form-group">
          <label>Email</label>
          <input type="email" className="form-control" name="email" onChange={this.handleChange} required />
        </div>
        <div className="form-group">
          <label>Password</label>
          <input type="password" className="form-control" name="password" onChange={this.handleChange} required />
        </div>
        <div className="form-group">
          <label>Confirm your password</label>
          <input type="password" className="form-control" name="password_confirm" onChange={this.handleChange} required />
        </div>
        <button type="submit" className="btn btn-primary">Submit</button>
      </form>
    );
  }
}

const mapDispatchToProps = dispatch => ({
  postSignup: body => dispatch(postSignup(body)),
});

export default connect(null, mapDispatchToProps)(Signup);

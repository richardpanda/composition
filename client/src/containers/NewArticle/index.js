import React, { Component } from 'react';
import { connect } from 'react-redux';

class NewArticle extends Component {
  constructor(props) {
    super(props);

    this.state = {
      title: '',
      body: '',
      success: '',
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

    this.setState({ success: '', error: '' });

    const { token } = this.props;
    const { title, body } = this.state;

    const requestBody = { title, body };
    const init = {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify(requestBody),
    };

    try {
      const response = await fetch('/api/articles', init);
      const responseBody = await response.json();

      if (response.ok) {
        this.setState({ success: 'Successfully added article!' });
      } else {
        this.setState({ error: responseBody.message });
      }
    } catch (e) {
      this.setState({ error: e });
    }
  }

  render() {
    const { success, error } = this.state;

    return (
      <form className="w-75 mx-auto" onSubmit={this.handleSubmit}>
        <h3 className="text-center">Creating New Article</h3>
        {success && <div className="alert alert-success">{success}</div>}
        {error && <div className="alert alert-danger">{error}</div>}
        <div className="form-group">
          <label>Title</label>
          <input type="text" className="form-control" name="title" onChange={this.handleChange} required />
        </div>
        <div className="form-group">
          <label>Body</label>
          <textarea className="form-control" name="body" rows="10" required onChange={this.handleChange}></textarea>
        </div>
        <button type="submit" className="btn btn-primary">Submit</button>
      </form>
    );
  }
}

const mapStateToProps = state => ({
  token: state.auth.token,
});

export default connect(mapStateToProps)(NewArticle);

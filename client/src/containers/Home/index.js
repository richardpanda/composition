import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';

import ArticlePreviews from '../ArticlePreviews';

class Home extends Component {
  render() {
    const { isLoggedIn } = this.props;

    const buttonStyle = {
      visibility: isLoggedIn ? '' : 'hidden',
    };

    return (
      <div className="pl-auto">
        <div className="row">
          <div className="col"></div>
          <h3 className="col text-center">Latest Articles</h3>
          <Link to="/articles/new" className="col text-right" style={buttonStyle}>
            <button className="btn btn-info btn-sm">Create Article</button>
          </Link>
        </div>
        <ArticlePreviews />
      </div>
    );
  }
}

const mapStateToProps = state => ({
  isLoggedIn: state.auth.isLoggedIn,
});

export default connect(mapStateToProps)(Home);

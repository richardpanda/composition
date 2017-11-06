import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';

class Home extends Component {
  render() {
    const { isLoggedIn } = this.props;

    return (
      <div className="pl-auto">
        {isLoggedIn &&
          <div className="d-flex flex-row-reverse">
            <Link to="/articles/new"><button className="btn btn-info btn-sm">Create Article</button></Link>
          </div>
        }
      </div>
    );
  }
}

const mapStateToProps = state => ({
  isLoggedIn: state.auth.isLoggedIn,
});

export default connect(mapStateToProps)(Home);

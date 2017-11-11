import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';

class ArticlePreviews extends Component {
  constructor(props) {
    super(props);

    this.state = {
      articlePreviews: [],
      error: ''
    };
    this.fetchArticlePreviews = this.fetchArticlePreviews.bind(this);
  }

  componentDidMount() {
    this.fetchArticlePreviews(this.props.page);
  }

  componentWillUpdate(nextProps, nextState) {
    if (this.props.page !== nextProps.page) {
      this.fetchArticlePreviews(nextProps.page);
    }
  }

  async fetchArticlePreviews(page) {
    try {
      const response = await fetch(`/api/articles?page=${page}`);
      const body = await response.json();

      if (response.ok) {
        this.setState({ articlePreviews: body.article_previews });
      } else {
        this.setState({ error: body.message });
      }
    } catch (e) {
      this.setState({ error: e });
    }
  }

  render() {
    const { articlePreviews, error } = this.state;
    const { page } = this.props;

    const showPrevious = page > 1;
    const showNext = articlePreviews.length === 10;

    return (
      <div>
        {error && <div className="alert alert-danger">{error}</div>}
        {articlePreviews.map(a => (
          <div className="card mb-2" key={a.article_id}>
            <div className="card-body p-2">
              <Link to={`/articles/${a.article_id}`}><h5 className="card-title mb-0">{a.title}</h5></Link>
              <p className="card-subtitle text-muted"><small>{a.username}</small></p>
            </div>
          </div>
        ))}
        <div>
          {showPrevious && <Link to={`?page=${page-1}`}><button className="btn btn-secondary btn-small mr-2">Previous</button></Link>}
          {showNext && <Link to={`?page=${page+1}`}><button className="btn btn-secondary btn-small">Next</button></Link>}
        </div>
      </div>
    );
  }
}

const parsePageNumber = (queryString) => {
  const pageRegex = /page=(\d+)/i;
  return (queryString.match(pageRegex) && parseInt(pageRegex.exec(queryString)[1], 10)) || 1;
};

const mapStateToProps = state => ({
  page: parsePageNumber(state.router.location.search),
});

export default connect(mapStateToProps)(ArticlePreviews);

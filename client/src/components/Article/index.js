import React, { Component } from 'react';

import './style.css';

class Article extends Component {
  constructor(props) {
    super(props);

    this.state = {
      article: null,
      error: '',
    };
  }

  async componentDidMount() {
    const { id } = this.props.match.params;

    try {
      const response = await fetch(`/api/articles/${id}`);
      const body = await response.json();

      if (response.ok) {
        this.setState({ article: body });
      } else {
        this.setState({ error: body.message });
      }
    } catch (e) {
      this.setState({ error: e });
    }
  }

  render() {
    const { article, error } = this.state;

    return (
      <div className="article-container mx-auto">
        {error && <div className="alert alert-danger">{error}</div>}
        {article &&
          <div>
            <div className="text-center mb-4">
              <h1 className="mb-0">{article.title}</h1>
              <h4><small className="text-muted">Written by: {article.username}</small></h4>
            </div>
            <div className="article-body">
              {article.body.split('\n').map((paragraph, idx) => (
                <div key={idx} className="mb-4">{paragraph}</div>
              ))}
            </div>
          </div>
        }
      </div>
    );
  }
}

export default Article;

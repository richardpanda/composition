import React, { Component } from 'react';

class ArticlePreviews extends Component {
  constructor(props) {
    super(props);

    this.state = {
      articlePreviews: [],
      error: ''
    };
  }

  async componentDidMount() {
    try {
      const response = await fetch('/api/articles');
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

    return (
      <div>
        {error && <div className="alert alert-danger">{error}</div>}
        {articlePreviews.map(a => (
          <div className="card mb-2" key={a.article_id}>
            <div className="card-body p-2">
              <h5 className="card-title mb-0">{a.title}</h5>
              <p className="card-subtitle text-muted"><small>{a.username}</small></p>
            </div>
          </div>
        ))}
      </div>
    );
  }
}

export default ArticlePreviews;

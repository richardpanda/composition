import React, { Component } from 'react';

import './App.css';
import ArticlePreviewList from './components/ArticlePreviewList';
import Navbar from './components/Navbar';

class App extends Component {
  render() {
    return (
      <div className="Container">
        <header>
          <Navbar />
        </header>
        <main>
          <h1 className="LatestArticles">Latest Articles</h1>
          <ArticlePreviewList />
        </main>
      </div>
    );
  }
}

export default App;

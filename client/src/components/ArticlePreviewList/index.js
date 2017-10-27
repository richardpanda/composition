import React, { Component } from 'react';

import './style.css';

class ArticlePreviewList extends Component {
  render() {
    const articlePreviews = [
      { username: "benwolford", title: "The product Facebook sells is you", article_id: 1, created_at: Date.now() },
      { username: "richardwhiteside", title: "Trapped in a web of algorithms", article_id: 2, created_at: Date.now() },
      { username: "thenewyorktimes", title: "Tech Giants Are Paying Huge Salaries for Scarce AI Talent", article_id: 3, created_at: Date.now() },
      { username: "bloomberg", title: "Elon Musk Was Wrong About Self-Driving Teslas", article_id: 4, created_at: Date.now() },
      { username: "meghanhebel", title: "Stop Sabotaging Your Code…Before You Even Code", article_id: 5, created_at: Date.now() },
      { username: "vaidehijoshi", title: "Less Repetition, More Dynamic Programming", article_id: 6, created_at: Date.now() },
      { username: "dianahsieh", title: "Learning w/ Diagrams: Handling Contention with PostgreSQL", article_id: 7, created_at: Date.now() },
      { username: "cezarykaraś", title: "Naming 101: Programmer’s Guide on How to Name Things", article_id: 8, created_at: Date.now() },
      { username: "gregsabo", title: "How awesome engineers ask for help", article_id: 9, created_at: Date.now() },
      { username: "kanishkdudeja", title: "The beauty of Go", article_id: 10, created_at: Date.now() },
    ];

    return (
      <ol className="ArticlePreviewList-list">
        {articlePreviews.map(({ username, title, article_id, created_at }) => (
          <li
            className="ArticlePreviewList-item"
            key={article_id}
          >
            <div className="ArticlePreviewList-title">{title}</div>
            <div className="ArticlePreviewList-username">@{username}</div>
            <div className="ArticlePreviewList-date">{new Date(created_at).toString()}</div>
          </li>
        ))}
      </ol>
    );
  }
}

export default ArticlePreviewList;

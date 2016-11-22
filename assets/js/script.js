'use strict';

import React from 'react';
import ReactDOM from 'react-dom';
import RealtimeVisitsCount from './components/realtime-visits.js';
import VisitsList from './components/visits-list.js';
import PageviewsList from './components/pageviews.js';
import VisitsGraph from './components/visits-graph.js';
import Login from './components/login.js';


class App extends React.Component {

  constructor(props) {
    super(props)
    this.state = { idToken: null }
  }

  render() {
    if(this.state.idToken) {
      return (
        <div className="container">
          <h1>Ana</h1>
          <RealtimeVisitsCount />
          <VisitsGraph />
          <PageviewsList />
        </div>
      );
    } else {
      return (
        <div className="container">
          <Login />
        </div>
      );
    }
  }
}

ReactDOM.render(
  <App />,
  document.getElementById('root')
);

'use strict';

import React from 'react';
import ReactDOM from 'react-dom';
import RealtimeVisitsCount from './components/realtime-visits.js';
import VisitsList from './components/visits-list.js';
import PageviewsList from './components/pageviews.js';
import VisitsGraph from './components/visits-graph.js';

ReactDOM.render(
  <div className="container">
    <h1>Ana</h1>
    <RealtimeVisitsCount />
    <VisitsGraph />
    <PageviewsList />
  </div>,
  document.getElementById('root')
);

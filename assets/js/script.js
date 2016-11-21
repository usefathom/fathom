'use strict';

import React from 'react';
import ReactDOM from 'react-dom';
import RealtimeVisitsCount from './components/realtime-visits.js';
import VisitsList from './components/visits-list.js';
import PageviewsList from './components/pageviews.js';

function tick() {
  const element = (
     <div className="container">
       <h1>Ana</h1>
       <RealtimeVisitsCount />
       <PageviewsList />
       <VisitsList />
     </div>
   );

  ReactDOM.render(
    element, document.getElementById('root')
  );
}

tick() && window.setInterval(tick, 1000);

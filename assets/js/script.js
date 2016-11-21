'use strict';

import React from 'react';
import ReactDOM from 'react-dom';
import RealtimeVisitsCount from './components/realtime-visits.js';
import VisitsList from './components/visits-list.js';

function tick() {
  const element = (
     <div className="container">
       <h1>Hello, world!</h1>
       <p>It is {new Date().toLocaleTimeString()}.</p>
       <RealtimeVisitsCount />
       <VisitsList />
     </div>
   );

  ReactDOM.render(
    element, document.getElementById('root')
  );
}

tick() && window.setInterval(tick, 1000);

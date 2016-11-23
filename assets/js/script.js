'use strict';

const m = require('mithril');
import Login from './components/login.js';
import Pageviews from './components/pageviews.js';
import RealtimeVisits from './components/realtime.js';
import VisitsGraph from './components/visits-graph.js';
import LogoutButton from './components/logoutButton.js';

const App = {
  controller(args) {
    this.state = {
      authenticated: false
    };

    this.setState = function(nextState) {
        m.startComputation();
        for(var k in nextState) {
          this.state[k] = nextState[k];
        }
        m.endComputation();
    }
  },
  view(c) {
    if( ! c.state.authenticated ) {
      return m.component(Login, {
          onAuth: () => {
            c.setState({ authenticated: true })
           }
        });
    }

    return [
      m('div.container', [
        m('div.header.cf', [
          m('h1.pull-left', 'Ana'),
          m('div.pull-right', [
            m.component(LogoutButton, {
              cb: () => {
                 c.setState({ authenticated: false })
               }
             })
          ]),
        ]),
        m.component(RealtimeVisits),
        m.component(VisitsGraph),
        m.component(Pageviews),
      ])
    ]
  }
}


m.mount(document.getElementById('root'), App)

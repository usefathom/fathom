'use strict';

import m from 'mithril';

function fetchData() {
  return fetch('/api/visits/count/realtime', {
    credentials: 'include'
  })
    .then((r) => r.json())
    .then((data) => {
      this.count = data;
      m.redraw();
  });
}

const RealtimeVisits = {
  controller(args) {
      this.count = 0;
      fetchData.bind(this) && window.setInterval(fetchData.bind(this), 6000);
  },

  view(c) {
    let visitors = c.count > 1 ? 'visitors' : 'visitor';

    return m('div.block', [
      m('span.count', c.count),
      ' ',
      m('span', visitors + " on the site right now.")
    ])
  }
}

export default RealtimeVisits

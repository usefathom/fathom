import m from 'mithril';

function fetchRecords() {
  return fetch('/api/pageviews', {
    credentials: 'include'
  }).then((r) => {
    if( r.ok ) {
      return r.json();
    }
  }).then((data) => {
    m.startComputation();
    this.records(data);
    m.endComputation();
  });
}

const Pageviews = {
  controller() {
    this.records = m.prop([]);
    fetchRecords.call(this) && window.setInterval(fetchRecords.bind(this), 60000);
  },
  view(c) {
    const tableRows = c.records().map((p, i) => m('tr', [
      m('td', i+1),
      m('td', [
        m('a', { href: p.Path }, p.Path)
      ]),
      m('td', p.Count),
      m('td', p.CountUnique)
    ]));

    return m('div.block', [
      m('h2', 'Pageviews'),
      m('table.table.pageviews', [
        m('thead', [
          m('tr', [
            m('th', '#'),
            m('th', 'URL'),
            m('th', 'Pageviews'),
            m('th', 'Unique'),
          ]) // tr
        ]), // thead
        m('tbody', tableRows )
      ]) // table
    ])
  }
}

export default Pageviews

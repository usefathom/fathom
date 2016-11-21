'use strict';

import React from 'react';
import ReactDOM from 'react-dom';


class VisitList extends React.Component {
  constructor(props) {
    super(props);
    this.state = { records: [] };
    this.refresh() && window.setInterval(this.refresh.bind(this), 1000);
  }

  refresh() {
    fetch('/api/visits')
      .then((r) => r.json())
      .then((data) => {
        this.setState({records: data});
    });
  }

  render() {
    const tableRows = this.state.records.map((visit) =>
      <tr key={visit.ID}>
        <td>{visit.Timestamp}</td>
        <td>{visit.IpAddress}</td>
        <td>{visit.Path}</td>
      </tr>
    );

    return (
      <div>
        <h2>Visits</h2>
        <table className="visits-table">
          <thead>
            <tr>
              <th>When</th>
              <th>IP Address</th>
              <th>Path</th>
            </tr>
          </thead>
          <tbody>{tableRows}</tbody>
        </table>
      </div>
    );
  }
}

function tick() {
  const element = (
     <div className="container">
       <h1>Hello, world!!</h1>
       <p>It is {new Date().toLocaleTimeString()}.</p>
       <VisitList />
     </div>
   );

  ReactDOM.render(
    element, document.getElementById('root')
  );
}

tick() && window.setInterval(tick, 1000);

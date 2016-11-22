import React, { Component } from 'react'

class VisitsList extends React.Component {
  constructor(props) {
    super(props);
    this.state = { records: [] };
    this.refresh() && window.setInterval(this.refresh.bind(this), 60000);
  }

  refresh() {
    return fetch('/api/visits')
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
      <div className="block">
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

export default VisitsList

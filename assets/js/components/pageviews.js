import React, { Component } from 'react'

class PageviewsList extends React.Component {
  constructor(props) {
    super(props);
    this.state = { records: [] };
    this.refresh() && window.setInterval(this.refresh.bind(this), 60000);
  }

  refresh() {
    return fetch('/api/pageviews')
      .then((r) => r.json())
      .then((data) => {
        this.setState({records: data});
    });
  }

  render() {
    const tableRows = this.state.records.map((p, i) =>
      <tr key={i}>
        <td>{i+1}</td>
        <td>{p.Path}</td>
        <td>{p.Count}</td>
      </tr>
    );

    return (
      <div>
        <h2>Pageviews</h2>
        <table className="table pageviews">
          <thead>
            <tr>
              <th>#</th>
              <th>Path</th>
              <th>Count</th>
            </tr>
          </thead>
          <tbody>{tableRows}</tbody>
        </table>
      </div>
    );
  }
}

export default PageviewsList

'use strict';

import { h, render, Component } from 'preact';

class Pageviews extends Component {

  constructor(props) {
    super(props)

    this.state = {
      records: []
    }
    this.fetchRecords = this.fetchRecords.bind(this);
    this.fetchRecords();
  }

  fetchRecords() {
    return fetch('/api/pageviews', {
      credentials: 'include'
    }).then((r) => {
      if( r.ok ) {
        return r.json();
      }
    }).then((data) => {
      this.setState({ records: data })
    });
  }

  render() {
    const tableRows = this.state.records.map( (p, i) => (
      <tr>
        <td>{i+1}</td>
        <td><a href={p.Path}>{p.Path}</a></td>
        <td>{p.Count}</td>
        <td>{p.CountUnique}</td>
      </tr>
    ));

    return (
      <div class="block">
        <h2>Pageviews</h2>
        <table class="table pageviews">
          <thead>
            <tr>
              <th>#</th>
              <th>URL</th>
              <th>Pageviews</th>
              <th>Unique</th>
            </tr>
          </thead>
          <tbody>{tableRows}</tbody>
        </table>
      </div>
    )
  }
}

export default Pageviews

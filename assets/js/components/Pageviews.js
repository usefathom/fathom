'use strict';

import { h, render, Component } from 'preact';

class Pageviews extends Component {

  constructor(props) {
    super(props)

    this.state = {
      records: []
    }
    this.fetchRecords = this.fetchRecords.bind(this);
    this.fetchRecords(props.period);
  }

  componentWillReceiveProps(newProps) {
    if(this.props.period != newProps.period) {
      this.fetchRecords(newProps.period)
    }
  }

  fetchRecords(period) {
    return fetch('/api/pageviews?period=' + period, {
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
        <td class="muted">{i+1}</td>
        <td><a href={p.Path}>{p.Path}</a></td>
        <td>{p.Count}</td>
        <td>{p.CountUnique}</td>
      </tr>
    ));

    return (
      <div class="block">
        <h3>Pageviews</h3>
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

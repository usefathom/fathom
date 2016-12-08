'use strict';

import { h, render, Component } from 'preact';
import * as numbers from '../lib/numbers.js';
import Client from '../lib/client.js';

const dayInSeconds = 60 * 60 * 24;

class Pageviews extends Component {

  constructor(props) {
    super(props)

    this.state = {
      records: []
    }
    this.fetchRecords = this.fetchRecords.bind(this);
  }

  componentDidMount() {
    this.fetchRecords(this.props.period);
  }

  componentWillReceiveProps(newProps) {
    if(this.props.period != newProps.period) {
      this.fetchRecords(newProps.period)
    }
  }

  fetchRecords(period) {
    const before = Math.round((+new Date() ) / 1000);
    const after = before - ( period * dayInSeconds );

    Client.request(`/pageviews?before=${before}&after=${after}`)
    .then((d) => { this.setState({ records: d })})
    .catch((e) => { console.log(e) })
  }

  render() {
    const tableRows = this.state.records.map( (p, i) => (
      <tr>
        <td class="muted">{i+1}</td>
        <td><a href={p.Path}>{p.Path.substring(0, 50)}{p.Path.length > 50 ? '..' : ''}</a></td>
        <td>{numbers.formatWithComma(p.Count)}</td>
        <td>{numbers.formatWithComma(p.CountUnique)}</td>
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

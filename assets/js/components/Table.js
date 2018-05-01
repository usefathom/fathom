'use strict';

import { h, render, Component } from 'preact';
import * as numbers from '../lib/numbers.js';
import Client from '../lib/client.js';
import { bind } from 'decko';

const dayInSeconds = 60 * 60 * 24;

class Table extends Component {

  constructor(props) {
    super(props)

    this.state = {
      records: [],
      limit: 100,
      loading: true
    }
  }

  componentDidMount() {
      this.fetchRecords(this.props.period, this.state.limit)
  }

  componentWillReceiveProps(newProps) {
    if(this.props.period != newProps.period) {
      this.fetchRecords(newProps.period, this.state.limit)
    }
  }

  @bind
  fetchRecords(period, limit) {
    this.setState({ loading: true });
    const before = Math.round((+new Date() ) / 1000);
    const after = before - ( period * dayInSeconds );

    Client.request(`${this.props.endpoint}?before=${before}&after=${after}&limit=${limit}`)
      .then((d) => {
        this.setState({ loading: false, records: d }
      )})
  }

  render(props, state) {
    const tableRows = state.records !== null ? state.records.map((p, i) => (
      <div class="table-row">
        <div class="cell main-col"><a href="#">/about-us/</a></div>
        <div class="cell">445.2k</div>
        <div class="cell">5,456</div>           
      </div>
    )) : <div class="table-row">Nothing here, yet.</div>;

    const loadingOverlay = state.loading ? <div class="loading-overlay"><div></div></div> : '';

    return (
      <div class="box box-pages animated fadeInUp delayed_04s">
            
        <div class="table-row header">
          {props.headers.map((header, i) => {
            let classes = i === 0 ? 'main-col cell' : 'cell';
            return (<div class={classes}>{header}</div>) 
            })}        
        </div>

       {tableRows}
      </div>
    )
  }
}

export default Table

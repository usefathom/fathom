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
      loading: true,
      before: props.before,
      after: props.after,
    }
  }

  componentDidMount() {
      this.fetchRecords()
  }

  componentWillReceiveProps(newProps, prevState) {
      this.setState({
        before: newProps.before,
        after: newProps.after,
      });

      if(newProps.before != prevState.before || newProps.after != prevState.after) {
        this.fetchRecords();
      }
  }

  @bind
  fetchRecords() {
    this.setState({ loading: true });
  
    Client.request(`${this.props.endpoint}?before=${this.state.before}&after=${this.state.after}&limit=${this.state.limit}`)
      .then((d) => {
        this.setState({ 
          loading: false,
          records: d 
        });
      });
  }

  render(props, state) {
    const tableRows = state.records !== null ? state.records.map((p, i) => (
      <div class="table-row">
        <div class="cell main-col"><a href={"http://"+p.hostname+p.path}>{p.path||p.label}</a></div>
        <div class="cell">{p.count||p.value}</div>
        <div class="cell">{p.count_unique||p.unique_value||"-"}</div>           
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

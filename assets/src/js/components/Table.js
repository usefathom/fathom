'use strict';

import { h, Component } from 'preact';
import * as numbers from '../lib/numbers.js';
import Client from '../lib/client.js';
import { bind } from 'decko';

const dayInSeconds = 60 * 60 * 24;

class Table extends Component {

  constructor(props) {
    super(props)

    this.state = {
      records: [],
      limit: 15,
      loading: true,
      total: 0,
    }
  }

  componentWillReceiveProps(newProps) {
      if(newProps.before == this.props.before && newProps.after == this.props.after) {
        return;
      }

      this.fetchRecords(newProps.before, newProps.after);
  }

  @bind
  fetchRecords(before, after) {
    this.setState({ loading: true });
  
    Client.request(`${this.props.endpoint}?before=${before}&after=${after}&limit=${this.state.limit}`)
      .then((d) => {
         // request finished; check if timestamp range is still the one user wants to see
        if( this.props.before != before || this.props.after != after ) {
          return;
        }

        this.setState({ 
          loading: false,
          records: d,
        });
      });

     // fetch totals too
     Client.request(`${this.props.endpoint}/pageviews?before=${before}&after=${after}`)
      .then((d) => {
        this.setState({ 
          total: d
        });
      });

  }

  render(props, state) {
    const tableRows = state.records !== null && state.records.length > 0 ? state.records.map((p, i) => {

      let href = (p.Hostname + p.Pathname) || p.URL;
      let classes = "table-row"; 
      if(state.total > 0) {
        classes += " w" + Math.min(98, Math.round(p.Pageviews / state.total * 100 * 2.5));
      }

      let label = p.Pathname
      if( props.showHostname ) {
        if( p.Group) {
          label = p.Group
        } else {
          label = p.Hostname.replace('www.', '').replace('https://', '').replace('http://', '') + (p.Pathname.length > 1 ? p.Pathname : '')
        }
      }

      return(
      <div class={classes}>
        <div class="cell main-col"><a href={href}>{label}</a></div>
        <div class="cell">{numbers.formatPretty(p.Pageviews)}</div>
        <div class="cell">{numbers.formatPretty(p.Visitors)||"-"}</div>           
      </div>
    )}) : <div class="table-row"><div class="cell main-col">Nothing here, yet.</div></div>;

    return (
      <div class={(state.loading ? "loading" : '')}>
        <div class="table-row header">
          {props.headers.map((header, i) => {
            return (<div class={i === 0 ? 'main-col cell' : 'cell'}>{header}</div>) 
            })}        
        </div>
        <div>
          {tableRows}
        </div>
      </div>
    )
  }
}

export default Table

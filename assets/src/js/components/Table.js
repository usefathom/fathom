'use strict';

import { h, Component } from 'preact';
import * as numbers from '../lib/numbers.js';
import Client from '../lib/client.js';
import { bind } from 'decko';
import classNames from 'classnames';
import { runInNewContext } from 'vm';

const dayInSeconds = 60 * 60 * 24;

class Table extends Component {

  constructor(props) {
    super(props)

    this.state = {
      records: [],
      offset: 0,
      limit: 15,
      loading: true,
      total: 0,
    }
  }

  componentWillReceiveProps(newProps, newState) {
    if(!this.paramsChanged(this.props, newProps)) {
      return;
    }

    this.fetchData(newProps)
  }

  paramsChanged(o, n) {
    return o.siteId != n.siteId || o.before != n.before || o.after != n.after;
  }
  
  @bind
  fetchData(props) {
    this.setState({ loading: true });

    Client.request(`/sites/${props.siteId}/stats/${props.endpoint}/agg?before=${props.before}&after=${props.after}&offset=${this.state.offset}&limit=${this.state.limit}`)
      .then((d) => {
         // request finished; check if timestamp range is still the one user wants to see
        if( this.paramsChanged(props, this.props) ) {
          return;
        }

        this.setState({
          loading: false,
          records: d,
        });
      });

     // fetch totals too
     Client.request(`/sites/${props.siteId}/stats/${props.endpoint}/agg/pageviews?before=${props.before}&after=${props.after}`)
      .then((d) => {
        this.setState({
          total: d
        });
      });
  }

  @bind 
  paginateNext() {
    this.setState({ offset: this.state.offset + this.state.limit })
    this.fetchData(this.props)
  }

  @bind 
  paginatePrev() {
    if(this.state.offset == 0) {
      return;
    }

    this.setState({ offset: Math.max(0, this.state.offset - this.state.limit) })
    this.fetchData(this.props)
  }

  render(props, state) {
    const tableRows = state.records !== null && state.records.length > 0 ? state.records.map((p, i) => {
      let href = (p.Hostname + p.Pathname) || p.URL;
      let widthClass = "";
      if(state.total > 0) {
        widthClass = "w" + Math.min(98, Math.round(p.Pageviews / state.total * 100 * 2.5));
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
      <div class={classNames("table-row", widthClass)}>
        <div class="cell main-col"><a href={href}>{label}</a></div>
        <div class="cell">{numbers.formatPretty(p.Pageviews)}</div>
        <div class="cell">{numbers.formatPretty(p.Visitors)||"-"}</div>
      </div>
    )}) : <div class="table-row"><div class="cell main-col">Nothing here, yet.</div></div>;

  // pagination row: only show when total # of results doesn't fit in one table page  
  const pagination = tableRows.length == state.limit || state.offset >= state.limit ? (
    <div class="row pag">
      <a href="javascript:void(0)" onClick={this.paginatePrev} class="back">‹</a>
      <a href="javascript:void(0)" onClick={this.paginateNext} class="next right">›</a>
    </div>) : '';

    return (
      <div class={classNames({ loading: state.loading })}>
        <div class="table-row header">
          {props.headers.map((header, i) => {
            return <div class={classNames("cell", { "main-col": i === 0 })}>{header}</div>
          })}
        </div>
        <div>
          {tableRows}
          {pagination}
        </div>
      </div>
    )
  }
}

export default Table

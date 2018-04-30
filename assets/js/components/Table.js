'use strict';

import { h, render, Component } from 'preact';
import * as numbers from '../lib/numbers.js';
import Client from '../lib/client.js';
const dayInSeconds = 60 * 60 * 24;

class Table extends Component {

  constructor(props) {
    super(props)

    this.state = {
      records: [],
      limit: 5,
      loading: true
    }

    this.tableHeaders = props.headers.map(heading => <th>{heading}</th>)
    this.fetchRecords = this.fetchRecords.bind(this)
    this.handleLimitChoice = this.handleLimitChoice.bind(this)
  }

  componentDidMount() {
      this.fetchRecords(this.props.period, this.state.limit)
  }

  labelCell(p) {
    if( this.props.labelCell ) {
      return this.props.labelCell(p)
    }

    return (
      <td>{p.Label.substring(0, 15)}</td>
    )
  }

  componentWillReceiveProps(newProps) {
    if(this.props.period != newProps.period) {
      this.fetchRecords(newProps.period, this.state.limit)
    }
  }

  handleLimitChoice(e) {
    this.setState({ limit: parseInt(e.target.value) })
    this.fetchRecords(this.props.period, this.state.limit)
  }

  fetchRecords(period, limit) {
    this.setState({ loading: true });
    const before = Math.round((+new Date() ) / 1000);
    const after = before - ( period * dayInSeconds );

    Client.request(`${this.props.endpoint}?before=${before}&after=${after}&limit=${limit}`)
      .then((d) => {
        this.setState({ loading: false, records: d }
      )}).catch((e) => { console.log(e) })
  }

  render() {
    const tableRows = this.state.records !== null ? this.state.records.map((p, i) => (
      <tr>
        <td class="muted">{i+1}</td>
        {this.labelCell(p)}
        <td>{numbers.formatWithComma(p.Value)}</td>
        <td>{Math.round(p.PercentageValue)}%</td>
      </tr>
    )) :<tr><td colspan="4" class="italic">Nothing here..</td></tr>;

    const loadingOverlay = this.state.loading ? <div class="loading-overlay"><div></div></div> : '';

    return (
      <div class="block">
        {loadingOverlay}
        <div class="clearfix">
          <h3 class="pull-left">{this.props.title}</h3>
          <div class="pull-right">
            <select onchange={this.handleLimitChoice}>
              <option>5</option>
              <option>20</option>
              <option>100</option>
            </select>
          </div>
          </div>
        <table>
          <thead>
            <tr>{this.tableHeaders}</tr>
          </thead>
          <tbody>
            {tableRows}
          </tbody>
        </table>
      </div>
    )
  }
}

export default Table

'use strict';

import { h, render, Component } from 'preact';
import * as numbers from '../lib/numbers.js';
const dayInSeconds = 60 * 60 * 24;

class Table extends Component {

  constructor(props) {
    super(props)

    this.state = {
      records: [],
      limit: 5
    }

    this.tableHeaders = props.headers.map(heading => <th>{heading}</th>)
    this.fetchRecords = this.fetchRecords.bind(this)
    this.handleLimitChoice = this.handleLimitChoice.bind(this)
    this.fetchRecords(props.period, this.state.limit)
  }

  labelCell(p) {
    if( this.props.labelCell ) {
      return this.props.labelCell(p)
    }

    return (
      <td>{p.Label}</td>
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
    const before = Math.round((+new Date() ) / 1000);
    const after = before - ( period * dayInSeconds );

    return fetch(`/api/${this.props.endpoint}?before=${before}&after=${after}&limit=${limit}`, {
      credentials: 'include'
    }).then((r) => {
      if( r.ok ) {
        return r.json();
      }

      // TODO: Make this pretty.
      if( r.status == 401 ) {
        this.props.onAuthError();
      }

      // TODO: do something with error
      throw new Error();
    }).then((data) => {
      this.setState({ records: data })
    }).catch((e) => {

    });
  }

  render() {
    const tableRows = this.state.records.map((p, i) => (
      <tr>
        <td class="muted">{i+1}</td>
        {this.labelCell(p)}
        <td>{numbers.formatWithComma(p.Count)}</td>
        <td>{Math.round(p.Percentage)}%</td>
      </tr>
    ));

    return (
      <div class="block">
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

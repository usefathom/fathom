'use strict';

import { h, render, Component } from 'preact';

class CountWidget extends Component {
  constructor(props) {
    super(props)

    this.state = {
      count: 0
    }

    this.fetchData = this.fetchData.bind(this);
    this.fetchData(props.period);
  }

  componentWillReceiveProps(newProps) {
    if(this.props.period != newProps.period) {
      this.fetchData(newProps.period)
    }
  }

  fetchData(period) {
    return fetch('/api/' + this.props.endpoint + '/count?period=' + period, {
      credentials: 'include'
    }).then((r) => {
        if( r.ok ) { return r.json(); }
        throw new Error();
     }).then((data) => {
        this.setState({ count: data })
    });
  }

  render() {
    return (

      <div class="block center-text">
        <h4 class="no-margin">{this.props.title}</h4>
        <div class="big">{this.state.count}</div>
        <div class="muted">last {this.props.period} days</div>
      </div>
    )
  }
}

export default CountWidget

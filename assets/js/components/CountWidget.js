'use strict';

import { h, render, Component } from 'preact';
const dayInSeconds = 60 * 60 * 24;

class CountWidget extends Component {
  constructor(props) {
    super(props)

    this.state = {
      count: 0,
      previousCount: 0
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
    const before = Math.round((+new Date() ) / 1000);
    const after = before - ( period * dayInSeconds );

    fetch(`/api/${this.props.endpoint}/count?before=${before}&after=${after}`, {
      credentials: 'include'
    }).then((r) => {
        if( r.ok ) { return r.json(); }
        throw new Error();
     }).then((data) => {
        this.setState({ count: data })
    });

    // query previous period
    const previousBefore = after;
    const previousAfter = previousBefore - ( period * dayInSeconds );
    fetch(`/api/${this.props.endpoint}/count?before=${previousBefore}&after=${previousAfter}`, {
      credentials: 'include'
    }).then((r) => {
        if( r.ok ) { return r.json(); }
        throw new Error();
     }).then((data) => {
        this.setState({ previousCount: data })
    });
  }

  renderPercentage() {
    if( ! this.state.previousCount ) {
      return '';
    }

    const percentage = Math.round(( this.state.count / this.state.previousCount * 100 - 100))
    return (
      <small class={percentage > 0 ? 'positive' : 'negative'}>{percentage}%</small>
    )
  }

  render() {
    return (
      <div class="block center-text">
        <h4 class="">{this.props.title}</h4>
        <div class="big tiny-margin">{this.state.count} {this.renderPercentage()}</div>
        <div class="muted">last {this.props.period} days</div>
      </div>
    )
  }
}

export default CountWidget

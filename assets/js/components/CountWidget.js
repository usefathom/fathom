'use strict';

import { h, Component } from 'preact';
import * as numbers from '../lib/numbers.js';
import Client from '../lib/client.js';
import { bind } from 'decko';

function getSundayOfCurrentWeek(d){
  var day = d.getDay();
  return new Date(d.getFullYear(), d.getMonth(), d.getDate() + (day == 0?0:7)-day );
}

const dayInSeconds = 60 * 60 * 24;

class CountWidget extends Component {
  constructor(props) {
    super(props)

    this.state = {
      value: '-',
      loading: false,
      before: props.before,
      after: props.after,
    }
  }

  componentDidMount() {
    this.fetchData();
  }

  componentWillReceiveProps(newProps, prevState) {
    if(newProps.before == prevState.before && newProps.after == prevState.after) {
      return;
    }

    this.setState({
      before: newProps.before,
      after: newProps.after,
      value: '-',
    });
    this.fetchData();
  }

  @bind
  fetchData() {
    this.setState({ loading: true })
    let before = this.state.before;
    let after = this.state.after;

    Client.request(`${this.props.endpoint}/count?before=${before}&after=${after}`)
      .then((d) => { 
        // request finished; check if timestamp range is still the one user wants to see
        if( this.state.before != before || this.state.after != after ) {
          return;
        }

        this.setState({ 
          loading: false, 
          value: numbers.formatWithComma(d), 
        })
      })
  }

  render(props, state) {
    return (
       <div class={"totals-detail " + ( state.loading ? "loading" : '')}>
        <div class="total-heading">{props.title}</div>
        <div class="total-numbers">{state.value}</div>
      </div>
    )
  }
}

export default CountWidget

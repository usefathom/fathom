'use strict';

import { h, render, Component } from 'preact';
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
      loading: false
    }
    this.fetchData(props.period);
  }

  componentWillReceiveProps(newProps) {
    console.log(newProps);
    if(this.props.period != newProps.period) {
      this.fetchData(newProps.period)
    }
  }

  @bind
  fetchData(period) {
    let before, after;
    let afterDate = new Date();
    afterDate.setHours(0, 0, 0, 0);
    switch(period) {
      case "week":
        afterDate.setDate(afterDate.getDate() - (afterDate.getDay() + 6) % 7);
      break;
      case "month":
        afterDate.setDate(1);
      break;
      case "year":
        afterDate.setDate(1);
        afterDate.setMonth(0);
      break;
    }

    before = Math.round((+new Date() ) / 1000);
    after = Math.round((+afterDate) / 1000);
    this.setState({ loading: true })

    Client.request(`${this.props.endpoint}/count?before=${before}&after=${after}`)
      .then((d) => { 
        this.setState({ 
          loading: false, 
          value: numbers.formatWithComma(d), 
        })
      })
  }

  render(props, state) {
    const loadingOverlay = state.loading ? <div class="loading-overlay"><div></div></div> : '';
    return (
       <div class="totals-detail">
        {loadingOverlay}
        <div class="total-heading">{props.title}</div>
        <div class="total-numbers">{state.value}</div>
      </div>
    )
  }
}

export default CountWidget

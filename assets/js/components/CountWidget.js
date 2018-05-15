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
      value: 0.00,
      loading: false,
      before: props.before,
      after: props.after,
    }
  }

  componentWillReceiveProps(newProps, prevState) {
    if(newProps.before == prevState.before && newProps.after == prevState.after) {
      return;
    }

    this.setState({
      before: newProps.before,
      after: newProps.after,
    });
    this.fetchData();
  }

  @bind 
  countUp(toValue) {
    const duration = 1000;
    const easeOutQuint = function (t) { return 1+(--t)*t*t*t*t };
    const setState = this.setState.bind(this);
    const startValue = this.state.value;
    const diff = toValue - startValue;
    let startTime;

    const tick = function(t) {
      if(!startTime) { startTime = t; }
      let progress = ( t - startTime ) / duration;
      let newValue = Math.round(startValue + (easeOutQuint(progress) * diff));
      setState({
        value: newValue,
      })

      if(progress < 1) {
        window.requestAnimationFrame(tick);
      }
    }

    window.requestAnimationFrame(tick);
  }

  @bind
  fetchData() {
    this.setState({ loading: true })
    let before = this.state.before;
    let after = this.state.after;

    Client.request(`${this.props.endpoint}?before=${before}&after=${after}`)
      .then((d) => { 
        // request finished; check if timestamp range is still the one user wants to see
        if( this.state.before != before || this.state.after != after ) {
          return;
        }

        this.setState({ 
          loading: false,
        })
        this.countUp(d);
      })
  }

  render(props, state) {
    let formattedValue = "-";

    if(state.value > 0) {
      switch(props.format) {
        case "percentage":
          formattedValue = numbers.formatPercentage(state.value)
        break;  

        default:
        case "number":
            formattedValue = numbers.formatWithComma(state.value)
        break;

        case "duration":
          formattedValue = numbers.formatDuration(state.value)
        break;
      }
    }

    return (
       <div class={"totals-detail " + ( state.loading ? "loading" : '')}>
        <div class="total-heading">{props.title}</div>
        <div class="total-numbers">{formattedValue}</div>
      </div>
    )
  }
}

export default CountWidget

'use strict';

import { h, Component } from 'preact';
import * as numbers from '../lib/numbers.js';
import Client from '../lib/client.js';
import { bind } from 'decko';


class CountWidget extends Component {
  constructor(props) {
    super(props)

    this.state = {
      value: "-",
      loading: false,
    }
  }

  componentWillReceiveProps(newProps, newState) {
    if(newProps.before == this.props.before && newProps.after == this.props.after) {
      return;
    }

    this.fetchData(newProps.before, newProps.after);
  }

  // TODO: Move to component of its own
  @bind 
  countUp(toValue) { 
    const duration = 1000;
    const easeOutQuint = function (t) { return 1+(--t)*t*t*t*t };
    const setState = this.setState.bind(this);
    const startValue = isFinite(this.state.value) ? this.state.value : 0;
    const diff = toValue - startValue;
    let startTime = performance.now();

    const tick = function(t) {
      let progress = Math.min(( t - startTime ) / duration, 1);
      let newValue = startValue + (easeOutQuint(progress) * diff);
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
  fetchData(before, after) {
    this.setState({ loading: true })

    Client.request(`${this.props.endpoint}?before=${before}&after=${after}`)
      .then((d) => { 
        // request finished; check if timestamp range is still the one user wants to see
        if( this.props.before != before || this.props.after != after ) {
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

    if(isFinite(state.value)) {
      switch(props.format) {
        case "percentage":
          formattedValue = numbers.formatPercentage(state.value)
        break;  

        default:
        case "number":
            formattedValue = numbers.formatWithComma(Math.round(state.value))
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

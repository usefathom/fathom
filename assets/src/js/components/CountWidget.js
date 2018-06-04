'use strict';

import { h, Component } from 'preact';
import * as numbers from '../lib/numbers.js';
import { bind } from 'decko';


class CountWidget extends Component {
  constructor(props) {
    super(props)

    this.state = {
      value: "-"
    }
  }

  componentWillReceiveProps(newProps, newState) {
    if(newProps.value == this.props.value) {
      return;
    }

    this.countUp(newProps.value);
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
       <div class={"totals-detail " + ( props.loading ? "loading" : '')}>
        <div class="total-heading">{props.title}</div>
        <div class="total-numbers">{formattedValue}</div>
      </div>
    )
  }
}

export default CountWidget

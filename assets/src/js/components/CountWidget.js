'use strict';

import { h, Component } from 'preact';
import * as numbers from '../lib/numbers.js';
import { bind } from 'decko';

const duration = 600;
const easeOutQuint = function (t) { return 1+(--t)*t*t*t*t };

class CountWidget extends Component {
  componentWillReceiveProps(newProps, newState) {
    if(newProps.value == this.props.value) {
      return;
    }

    this.countUp(this.props.value || 0, newProps.value);
  }

  // TODO: Move to component of its own
  @bind 
  countUp(fromValue, toValue) { 
    const format = this.formatValue.bind(this);
    const startValue = isFinite(fromValue) ? fromValue : 0;
    const numberEl = this.numberEl;
    const diff = toValue - startValue;
    let startTime = performance.now();

    const tick = function(t) {
      let progress = Math.min(( t - startTime ) / duration, 1);
      let newValue = startValue + (easeOutQuint(progress) * diff);
      numberEl.textContent = format(newValue)

      if(progress < 1) {
        window.requestAnimationFrame(tick);
      }
    }
    
    window.requestAnimationFrame(tick);
  }

  @bind
  formatValue(value) {
    let formattedValue = "-";

    if(isFinite(value)) {
      switch(this.props.format) {
        case "percentage":
          formattedValue = numbers.formatPercentage(value)
        break;  

        default:
        case "number":
          formattedValue = numbers.formatPretty(Math.round(value))
        break;

        case "duration":
          formattedValue = numbers.formatDuration(value)
        break;
      }
    }

    return formattedValue;
  }

  render(props, state) {
    return (
       <div class={"totals-detail " + ( props.loading ? "loading" : '')}>
        <div class="total-heading">{props.title}</div>
        <div class="total-numbers" ref={(e) => { this.numberEl = e; }}>{this.formatValue(props.value)}</div>
      </div>
    )
  }
}

export default CountWidget

'use strict';

import { h, render, Component } from 'preact';
import { bind } from 'decko';

const availablePeriods = [
  {
    id: 'day',
    label: 'Today'
  },
  {
    id: 'week',
    label: 'This week'
  },
  {
    id: 'month',
    label: 'This month'
  },
  {
    id: 'year',
    label: 'This year'
  }
]

class DatePicker extends Component {
  constructor(props) {
    super(props)

    this.state = {
      period: this.props.value
    }
  }

  @bind
  setPeriod(e) {
    e.preventDefault();

    var nextState = { 
      period: e.target.dataset.value 
    }

    if(this.state.period != nextState.period) {
      this.setState(nextState)
      this.props.onChange(nextState.period);
    }
  }

  render(props, state) {
    const links = availablePeriods.map((p) => {
      let className = ( p.id == state.period ) ? 'active' : '';
      return <li class={className} ><a href="#" data-value={p.id} onClick={this.setPeriod}>{p.label}</a></li>
    });

    return (

      <ul>
        {links}
      </ul>
    )
  }
}

export default DatePicker

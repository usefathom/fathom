'use strict';

import { h, Component } from 'preact';
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
      period: props.value,
      before: 0,
      after: 0,
    }
  }

  componentDidMount() {
    this.setTimeRange(this.state.period)
  }

  @bind
  setTimeRange(period) {
    let beforeDate = new Date();
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

    let before, after;
    before = Math.round((+beforeDate ) / 1000);
    after = Math.round((+afterDate) / 1000);
    this.setState({
      period: period,
      before: before, 
      after: after,
    });
    this.props.onChange(this.state);
  }

  @bind
  setPeriod(e) {
    e.preventDefault();

    let newPeriod = e.target.dataset.value;
    if( newPeriod === this.state.period) {
      return;
    }

    this.setTimeRange(newPeriod);
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

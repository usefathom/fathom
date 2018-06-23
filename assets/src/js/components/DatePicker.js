'use strict';

import { h, Component } from 'preact';
import { bind } from 'decko';
import Pikadayer from './Pikadayer.js';

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
];

const padZero = function(n){return n<10? '0'+n:''+n;}

class DatePicker extends Component {
  constructor(props) {
    super(props)

    this.state = {
      period: props.value,
      before: 0,
      after: 0,
      startDate: null,
      endDate: null,
    }

    this.updateDatesFromPeriod(this.state.period)
  }

  @bind
  updateDatesFromPeriod(period) {
    let now = new Date();
    let startDate, endDate;

    switch(period) {
      case "day":
        startDate = new Date(now.getFullYear(), now.getMonth(), now.getDate());
        endDate = new Date(now.getFullYear(), now.getMonth(), now.getDate());
      break;

      case "week":
        startDate = new Date(now.getFullYear(), now.getMonth(), (now.getDate() - (now.getDay() + 6) % 7));
        endDate = new Date(startDate.getFullYear(), startDate.getMonth(), startDate.getDate() + 6);
      break;
      case "month":
        startDate = new Date(now.getFullYear(), now.getMonth(), 1);
        endDate = new Date(startDate.getFullYear(), startDate.getMonth() + 1, 0);
      break;
      case "year":
        startDate = new Date(now.getFullYear(), 0, 1);
        endDate = new Date(startDate.getFullYear() + 1, 0, 0);
      break;

      case "":
      default:
        return; // empty period, don't update state
      break;
    }

    startDate.setHours(0, 0, 0);
    endDate.setHours(23, 59, 59);
    this.setDateRange(startDate, endDate, period);
  }

  @bind
  setDateRange(startDate, endDate, period) {
    // don't update state if start > end. user may be busy picking dates.
    // todo: show error
    if(startDate > endDate) {
      return;
    }

    const timezoneOffset = (new Date()).getTimezoneOffset() * 60;

    let before, after;
    before = Math.round(((+endDate) / 1000) - timezoneOffset);
    after = Math.round(((+startDate) / 1000) - timezoneOffset);

    this.setState({
      period: period || '',
      startDate: startDate,
      endDate: endDate,
      before: before, 
      after: after,
    });

    // use slight delay for updating rest of application to allow this function to be called again
    if(!this.timeout) {
      this.timeout = window.setTimeout(() => {
        this.props.onChange(this.state);
        this.timeout = null;
      }, 2)
    }
  }

  @bind
  setPeriod(e) {
    e.preventDefault();

    let newPeriod = e.target.getAttribute('data-value');
    if( newPeriod === this.state.period) {
      return;
    }

    this.updateDatesFromPeriod(newPeriod);
  }

  dateValue(date) {
    return date.getFullYear() + '-' + padZero(date.getMonth() + 1) + '-' + padZero(date.getDate());
  }

  @bind 
  setStartDate(date) {
    if(date) {
      this.setDateRange(date, this.state.endDate, '')
    }
  }

  @bind 
  setEndDate(date) {
    if(date) {
      this.setDateRange(this.state.startDate, date, '')
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
        <li class="custom">
          <Pikadayer value={this.dateValue(state.startDate)} onSelect={this.setStartDate} />
          <span style="margin: 0 8px"> to </span> 
          <Pikadayer value={this.dateValue(state.endDate)} onSelect={this.setEndDate}  />
        </li>
      </ul>
    )
  }
}

export default DatePicker

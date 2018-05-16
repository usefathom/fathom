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

    this.updateDatesFromPeriod(this.state.period)
  }

  @bind
  updateDatesFromPeriod(period) {
    let startDate = new Date();
    startDate.setHours(0);
    startDate.setMinutes(0);

    let endDate = new Date();
    endDate.setHours(24);
    endDate.setMinutes(0);

    switch(period) {
      case "week":
        startDate.setDate(startDate.getDate() - (startDate.getDay() + 6) % 7);
        endDate.setDate(startDate.getDate() + 7);
      break;
      case "month":
        startDate.setDate(1);
        endDate.setMonth(endDate.getMonth() + 1);
        endDate.setDate(0);
      break;
      case "year":
        startDate.setDate(1);
        startDate.setMonth(0);
        endDate.setYear(startDate.getFullYear() + 1);
        endDate.setMonth(0);
        endDate.setDate(0);
      break;
    }

    this.setDateRange(startDate, endDate, period);
  }

  @bind
  setDateRange(startDate, endDate, period) {
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
      picking: '',
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

    this.updateDatesFromPeriod(newPeriod);
  }

  dateValue(date) {
    const addZero = function(n){return n<10? '0'+n:''+n;}
    return date.getFullYear() + '-' + addZero(date.getMonth() + 1) + '-' + addZero(date.getDate());
  }

  @bind
  startPicking(e) {
    this.setState({ picking: e.target.dataset.value })
  }

  @bind
  stopPicking(e) {
    this.setState({ picking: '' })
  }

  @bind 
  setStartDate(e) {
    let newStartDate = e.target.valueAsDate;
    if(newStartDate) {
      this.setDateRange(newStartDate, this.state.endDate, '')
    }
  }

  @bind 
  setEndDate(e) {
    let newEndDate = e.target.valueAsDate;
    if(newEndDate) {
      this.setDateRange(this.state.startDate, newEndDate, '')
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
        <li>
          <span style="padding: 0 8px 0 0;">&mdash;</span> 
          <span class="datepicker-wrap">
            <strong onclick={this.startPicking} data-value="start">{state.startDate.toLocaleDateString()}</strong>
            <span class="datepicker" style={state.picking === 'start' ? '' : 'display: none'}>
              <label>Choose start date</label>
              <input type="date" value={this.dateValue(state.startDate)} onblur={this.stopPicking} onchange={this.setStartDate} />
            </span>
          </span>
          <span> to </span> 
          <span class="datepicker-wrap">
            <strong onclick={this.startPicking} data-value="end">{state.endDate.toLocaleDateString()}</strong>
            <span class="datepicker" style={state.picking === 'end' ? '' : 'display: none'}>
              <label>Choose end date</label>
              <input type="date" value={this.dateValue(state.endDate)} onblur={this.stopPicking} onchange={this.setStartDate} />
            </span>
          </span>
        </li>
      </ul>
    )
  }
}

export default DatePicker

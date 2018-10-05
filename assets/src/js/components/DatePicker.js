'use strict';

import { h, Component } from 'preact';
import { bind } from 'decko';
import Pikadayer from './Pikadayer.js';
import classNames from 'classnames';

const defaultPeriod = 'last-7-days';
const padZero = function(n){return n<10? '0'+n:''+n;}

function getNow() {
  let now = new Date()
  let tzOffset = now.getTimezoneOffset()  * 60 * 1000;

  // if we're ahead of UTC, stick to UTC's "now"
  // this is ugly but a sad necessity for now because we store and aggregate statistics using UTC dates (without time data)
  // For those ahead of UTC, "today" will always be empty if they're checking early on in their day
  // see https://github.com/usefathom/fathom/issues/134
  if (tzOffset < 0) {
    now.setTime(now.getTime() + tzOffset )
  }

  return now
}

// today, yesterday, this week, last 7 days, last 30 days
const availablePeriods = {
  'today': {
    label: 'Today',
    start: function() {
      const now = getNow();
      return new Date(now.getFullYear(), now.getMonth(), now.getDate());
    },
    end: function() {
      const now = getNow();
      return new Date(now.getFullYear(), now.getMonth(), now.getDate());
    },
 },
  'last-7-days': {
    label: 'Last 7 days',
    start: function() {
      const now = getNow();
      return new Date(now.getFullYear(), now.getMonth(), now.getDate()-6);
    },
    end: function() {
      const now = getNow();
      return new Date(now.getFullYear(), now.getMonth(), now.getDate());
    },
 },
  'last-30-days': {
    label: 'Last 30 days',
    start: function() {
      const now = getNow();
      return new Date(now.getFullYear(), now.getMonth(), now.getDate()-29);
    },
    end: function() {
      const now = getNow();
      return new Date(now.getFullYear(), now.getMonth(), now.getDate());
    },
 },
  'this-year': {
    label: 'This year',
    start: function() {
      const now = getNow();
      return new Date(now.getFullYear(), 0, 1);
    },
    end: function() {
      const now = getNow();
      return new Date(this.start().getFullYear() + 1, 0, 0);
    },
 },
}

class DatePicker extends Component {
  constructor(props) {
    super(props)

    this.state = {
      period: window.location.hash.substring(2) || window.localStorage.getItem('period') || defaultPeriod,
      before: 0, // UTC timestamp
      after: 0, // UTC timestamp
      startDate: null, // local date object
      endDate: null, // local date object
    }

    this.updateDatesFromPeriod(this.state.period)
    window.addEventListener('keydown', this.handleKeyPress);
  }

  @bind
  updateDatesFromPeriod(period) {
    if(typeof(availablePeriods[period]) !== "object") {
      return;
    }
    let p = availablePeriods[period];
    this.setDateRange(p.start(), p.end(), period);
  }

  @bind
  setDateRange(startDate, endDate, period) {
    // don't update state if start > end. user may be busy picking dates.
    // todo: show error
    if(startDate > endDate) {
      return;
    }

    // include start & end day by forcing time
    startDate.setHours(0, 0, 0);
    endDate.setHours(23, 59, 59);

    // create unix timestamps from local date objects
    let before, after;
    const timezoneOffset = (new Date()).getTimezoneOffset() * 60;
    before = Math.round((+endDate) / 1000) - timezoneOffset;
    after = Math.round((+startDate) / 1000) - timezoneOffset;

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

        window.localStorage.setItem('period', this.state.period)
        window.history.replaceState(this.state, null, `#!${this.state.period}`)
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
    this.setDateRange(date, this.state.endDate, '')
  }

  @bind
  setEndDate(date) {
    this.setDateRange(this.state.startDate, date, '')
  }

  @bind
  handleKeyPress(evt) {
    // Don't handle input when the user is in a text field or text area.
    let tag = event.target.tagName;
    if(tag === "INPUT" || tag === "TEXTAREA") {
      return;
    }

    // TODO: Account for leap years
    let diff = this.state.endDate - this.state.startDate + 1000;
    let newStartDate, newEndDate;

    switch(evt.keyCode) {
      // left-arrow
      case 37:
        newStartDate = new Date(+this.state.startDate - diff)
        newEndDate = new Date(+this.state.endDate - diff)
        this.setDateRange(newStartDate, newEndDate)
      break;

      //right-arrow
      case 39:
      newStartDate = new Date(+this.state.startDate + diff)
      newEndDate = new Date(+this.state.endDate + diff)
      this.setDateRange(newStartDate, newEndDate)
      break;
    }
  }

  render(props, state) {
    const links = Object.keys(availablePeriods).map((id) => {
      let p = availablePeriods[id];
      return (
        <li class={classNames({ active: id == state.period })}>
          <a href="#" data-value={id} onClick={this.setPeriod}>{p.label}</a>
        </li>
      );
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

'use strict';

import { h, render, Component } from 'preact';

var availablePeriods = [
  {
    id: 7,
    label: 'Last 7 days'
  },
  {
    id: 30,
    label: 'Last 30 days'
  },
  {
    id: 90,
    label: 'Last quarter'
  },
  {
    id: 360,
    label: 'Last year'
  }
]

class DatePicker extends Component {
  constructor(props) {
    super(props)

    this.state = {
      period: this.props.period
    }

    this.setPeriod = this.setPeriod.bind(this)
  }

  setPeriod(e) {
    var nextState = { period: parseInt(e.target.value) }
    if(this.state.period != nextState.period) {
      this.setState(nextState)
      this.props.onChoose(this.state.period);
    }
  }

  render() {
    const buttons = availablePeriods.map((p) => {
      let className = ( p.id == this.state.period ) ? 'active' : '';
      return <button value={p.id} class={className} onClick={this.setPeriod}>{p.label}</button>
    });

    return (

      <div class="small-margin">
        {buttons}
      </div>
    )
  }
}

export default DatePicker

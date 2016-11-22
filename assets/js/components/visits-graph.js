'use strict';

import React, { Component } from 'react'
import Chart from 'chart.js'

Chart.defaults.global.tooltips.xPadding = 12;
Chart.defaults.global.tooltips.yPadding = 12;
//Chart.defaults.global.scales.yAxes.ticks.beginAtZero = true; // = 12;


class VisitsGraph extends React.Component {
  constructor(props) {
    super(props);
    this.state = { data: [] }
    this.refresh();
  }

  refresh() {
    return fetch('/api/visits/count/day')
      .then((r) => r.json())
      .then((data) => {
        this.setState({ data: data });
        this.initChart(this.chartCtx);
    });
  }

  initChart(ctx) {
     this.chart = new Chart(ctx, {
      type: 'line',
      data: {
        labels: this.state.data.map((d) => d.Label),
        datasets: [{
            label: '# of Visitors',
            data: this.state.data.map((d) => d.Count),
            backgroundColor: 'rgba(0, 155, 255, .2)'
        }]
      },
      options: {
        scale: {
          ticks: {
            beginAtZero: true
          }
        }
      }
    });
  }

  render() {
    return (
      <div className="block">
        <canvas ref={(ctx) => { this.chartCtx = ctx; }} width="600" height="200"></canvas>
      </div>
    );
  }
}

export default VisitsGraph;

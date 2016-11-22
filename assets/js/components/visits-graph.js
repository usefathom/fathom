'use strict';

import React, { Component } from 'react'
import Chart from 'chart.js'

function rand(min, max, num) {
    var rtn = [];
    while (rtn.length < num) {
      rtn.push((Math.random() * (max - min)) + min);
    }
    return rtn;
  }

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
    if( ! this.state.data.length ) {
      return;
    }

    console.log(this.state.data.map((d) => d.Label));
    console.log(this.state.data.map((d) => d.Count ));
    console.log("Init chart");
    console.log(ctx);

     var myChart = new Chart(ctx, {
      type: 'line',
      data: {
        labels: this.state.data.map((d) => d.Label),
        datasets: [{
            label: '# of Visitors',
            data: this.state.data.map((d) => d.Count),
            backgroundColor: 'rgba(162, 52, 235, 0.2)'
        }]
      },
      options: {
        tooltips: {
          xPadding: 12,
          yPadding: 9,
        },
        scales: {
            yAxes: [{
                ticks: {
                    beginAtZero:true
                }
            }]
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

'use strict';

import { h, render, Component } from 'preact';
import Chart from 'chart.js'

class Graph extends Component {

  constructor(props) {
    super(props)

    this.state = {
      data: []
    }
    this.fetchData = this.fetchData.bind(this);
    this.fetchData();
  }

  initChart() {
     new Chart(this.ctx, {
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

  fetchData() {
    return fetch('/api/visits/count/day', {
      credentials: 'include'
    })
      .then((r) => r.json())
      .then((data) => {
        this.setState({ data: data})
        this.initChart()
    });
  }

  render() {
    return (
      <div class="block">
        <canvas width="600" height="220" ref={(el) => { this.ctx = el; }} />
      </div>
    )
  }
}

export default Graph

'use strict';

import { h, render, Component } from 'preact';
import Chart from 'chart.js'

Chart.defaults.global.tooltips.xPadding = 10;
Chart.defaults.global.tooltips.yPadding = 10;

class Graph extends Component {

  constructor(props) {
    super(props)

    this.state = {
      visitorData: [],
      pageviewData: []
    }
    this.fetchData = this.fetchData.bind(this);
    this.fetchData();
  }

  refreshChart() {
     this.chart = new Chart(this.ctx, {
      type: 'line',
      data: {
        labels: this.state.visitorData.map((d) => d.Label),
        datasets: [
          {
            label: '# of Visitors',
            data: this.state.visitorData.map((d) => d.Count),
            backgroundColor: 'rgba(0, 0, 255, .2)'
          },
          {
            label: '# of Pageviews',
            data: this.state.pageviewData.map((d) => d.Count),
            backgroundColor: 'rgba(0, 0, 125, .2)'
          }
      ]
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
    // fetch visitor data
    fetch('/api/visits/count/day', {
      credentials: 'include'
    }).then((r) => r.json())
      .then((data) => {
        this.setState({ visitorData: data})
        this.refreshChart();
    });

    // fetch pageview data
    fetch('/api/pageviews/count/day', {
      credentials: 'include'
    }).then((r) => r.json())
      .then((data) => {
        this.setState({ pageviewData: data})
        this.refreshChart();
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

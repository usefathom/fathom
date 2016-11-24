'use strict';

import { h, render, Component } from 'preact';
import Chart from 'chart.js'

Chart.defaults.global.tooltips.xPadding = 10;
Chart.defaults.global.tooltips.yPadding = 10;
Chart.defaults.global.layout = { padding: 10 }

class Graph extends Component {
  constructor(props) {
    super(props)

    this.state = {
      visitorData: [],
      pageviewData: []
    }
    this.fetchData = this.fetchData.bind(this);
    this.fetchData(props.period);
  }

  componentWillReceiveProps(newProps) {
    if(this.props.period != newProps.period) {
      this.fetchData(newProps.period)
    }
  }

  refreshChart() {
    if( ! this.canvas ) { return; }
    if( this.chart ) { this.chart.clear(); }

    // clear canvas
    var newCanvas = document.createElement('canvas');
    newCanvas.setAttribute('width', this.canvas.getAttribute('width'));
    newCanvas.setAttribute('height', this.canvas.getAttribute('height'));
    this.canvas.parentNode.replaceChild(newCanvas, this.canvas);
    this.canvas = newCanvas;

    this.chart = new Chart(this.canvas, {
      type: 'line',
      data: {
        labels: this.state.visitorData.map((d) => d.Label),
        datasets: [
          {
            label: '# of Visitors',
            data: this.state.visitorData.map((d) => d.Count),
            backgroundColor: 'rgba(255, 155, 0, .6)',
            pointStyle: 'rect',
            pointBorderWidth: 0.1,
          },
          {
            label: '# of Pageviews',
            data: this.state.pageviewData.map((d) => d.Count),
            backgroundColor: 'rgba(0, 155, 255, .4)',
            pointStyle: 'rect',
            pointBorderWidth: 0.1,
          }
      ],
      }
    });
  }

  fetchData(period) {
    // fetch visitor data
    fetch('/api/visits/count/day?period=' + period, {
      credentials: 'include'
    }).then((r) => r.json())
      .then((data) => {
        this.setState({ visitorData: data })
        window.setTimeout(() => (this.refreshChart()), 20);
    });

    // fetch pageview data
    fetch('/api/pageviews/count/day?period=' + period, {
      credentials: 'include'
    }).then((r) => r.json())
      .then((data) => {
        this.setState({ pageviewData: data })
        window.setTimeout(() => (this.refreshChart()), 20);
    });
  }

  render() {
    return (
      <div class="block">
        <canvas ref={(el) => { this.canvas = el; }} />
      </div>
    )
  }
}

export default Graph

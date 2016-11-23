'use strict'

import m from 'mithril';
import Chart from 'chart.js'

Chart.defaults.global.tooltips.xPadding = 12;
Chart.defaults.global.tooltips.yPadding = 12;

function fetchData() {
  return fetch('/api/visits/count/day', {
    credentials: 'include'
  })
    .then((r) => r.json())
    .then((data) => {
      this.data = data;
      initChart.call(this);
  });
}

function initChart() {
  let ctx = this.chartCtx;

   new Chart(ctx, {
    type: 'line',
    data: {
      labels: this.data.map((d) => d.Label),
      datasets: [{
          label: '# of Visitors',
          data: this.data.map((d) => d.Count),
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

const VisitsGraph = {
  controller(args) {
    this.data = [];
    fetchData.call(this);
  },

  view(c) {
      return m('div.block', [
        m('canvas', {
          width: 600,
          height: 200,
          config: (el) => { c.chartCtx = el; }
        })
      ])
  },

}

export default VisitsGraph;

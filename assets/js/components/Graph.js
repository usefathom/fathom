'use strict';

import { h, render, Component } from 'preact';
import * as d3 from 'd3';
import tip from 'd3-tip';
d3.tip = tip;

const dayInSeconds = 60 * 60 * 24;

function Chart(element) {
  var padt = 10, padb = 20, padr = 40, padl = 40,
      h = 300,
      w = element.parentNode.clientWidth - padl - padr,
      x  = d3.scaleBand().range([0, w]).padding(0.2).round(true),
      y  = d3.scaleLinear().range([h, 0]),
      yAxis = d3.axisLeft().scale(y).tickSize(-w + padl + padr),
      xAxis = d3.axisBottom().scale(x);

  var pageviewTip = d3.tip()
      .attr('class', 'd3-tip')
      .html((d) => '<span>' + d.Count + '</span>' + ' pageviews')
      .offset([-12, 0]);

  var visitorTip = d3.tip()
      .attr('class', 'd3-tip')
      .html((d) => '<span>' + d.Count + '</span>' + ' visitors' )
      .offset([-12, 0]);

  var graph = d3.select('#graph');
  var vis = graph
      .append('svg')
      .attr('width', w + padl + padr)
      .attr('height', h + padt + padb)
      .append('g')
      .attr('transform', 'translate(' + padl + ',' + padt + ')');

  vis.call(pageviewTip);
  vis.call(visitorTip);

  function update(primaryData, secondaryData) {
    var max = d3.max(primaryData, (d) => d.Count);
    var ticks = primaryData.length;
    var xTick = Math.round(ticks / 7);

    x.domain(primaryData.map((d) => d.Label))
    y.domain([0, (max * 1.1)])

    // clear all previous data
    vis.selectAll('*').remove();

    // axes
    vis.append("g")
      .attr("class", "y axis")
      .call(yAxis);

    vis.append("g")
      .attr("class", "x axis")
      .attr('transform', 'translate(0,' + h + ')')
      .call(xAxis)
      .selectAll('g')
      .style('display', (d, i) => i % xTick != 0 ? 'none' : 'block')

    // bars
    var bars = vis.selectAll('g.primary-bar')
      .data(primaryData)
      .enter()
      .append('g')
      .attr('class', 'primary-bar')
      .attr('transform', function (d, i) { return "translate(" + x(d.Label) + ", 0)" });

    bars.append('rect')
      .attr('width', x.bandwidth())
      .attr('height', (d) => (h - y(d.Count)) )
      .attr('y', (d) => y(d.Count))
      .on('mouseover', pageviewTip.show)
      .on('mouseout', pageviewTip.hide);

    var visitorBars = vis.selectAll('g.sub-bar')
      .data(secondaryData)
      .enter()
      .append('g')
      .attr('class', 'sub-bar')
      .attr('transform', (d, i) => "translate(" + ( x(d.Label) + ( x.bandwidth() * 0.16667 ) ) + ", 0)");

    visitorBars.append('rect')
      .attr('width', x.bandwidth() * 0.66 )
      .attr('height', (d) => (h - y(d.Count)) )
      .attr('y', (d) => y(d.Count))
      .on('mouseover', visitorTip.show)
      .on('mouseout', visitorTip.hide);
  }

  return { 'update': update }
}


class Graph extends Component {
  constructor(props) {
    super(props)

    this.fetchData = this.fetchData.bind(this);
    this.refreshChart = this.refreshChart.bind(this);
    this.data = {
      visitors: null,
      pageviews: null,
    }
  }

  componentDidMount() {
    this.chart = new Chart(document.getElementById('graph'));
    this.fetchData(this.props.period);
  }

  componentWillReceiveProps(newProps) {
    if(this.props.period != newProps.period) {
      this.fetchData(newProps.period)
    }
  }

  shouldComponentUpdate() { return false; }

  refreshChart() {
    if(this.data.visitors && this.data.pageviews) {
      this.chart.update(this.data.pageviews, this.data.visitors);
    }
  }

  fetchData(period) {
    const before = Math.round((+new Date() ) / 1000);
    const after = before - ( period * dayInSeconds );
    const group = period > 90 ? 'month' : 'day';

    // fetch visitor data
    fetch(`/api/visits/count/group/${group}?before=${before}&after=${after}`, {
      credentials: 'include'
    }).then((r) => {
      if( r.ok ) {
        return r.json();
      }
      throw new Error();
    }).then((data) => {
      this.data.visitors = data;
      window.requestAnimationFrame(this.refreshChart);
    });

    // fetch pageview data
    fetch(`/api/pageviews/count/group/${group}?before=${before}&after=${after}`, {
      credentials: 'include'
    }).then((r) => {
      if( r.ok ) {
        return r.json();
      }

      throw new Error();
    }).then((data) => {
        this.data.pageviews = data;
        window.requestAnimationFrame(this.refreshChart);
    });
  }

  render() {
    return (
      <div class="block">
        <div id="graph"></div>
      </div>
    )
  }
}

export default Graph

'use strict';

import { h, render, Component } from 'preact';
import * as d3 from 'd3';
import tip from 'd3-tip';

d3.tip = tip;

const dayInSeconds = 60 * 60 * 24;

class Graph extends Component {
  constructor(props) {
    super(props)

    this.state = {
      visitorData: [],
      pageviewData: []
    }

    this.fetchData = this.fetchData.bind(this);
    this.refreshChart = this.refreshChart.bind(this);
  }

  componentDidMount() {
    this.fetchData(this.props.period);
  }

  componentWillReceiveProps(newProps) {
    if(this.props.period != newProps.period) {
      this.fetchData(newProps.period)
    }
  }

  refreshChart() {
    var padt = 10, padb = 20, padr = 40, padl = 40,
        h = 300,
        w = document.getElementById('graph').parentNode.clientWidth - padl-padr,
        x  = d3.scaleBand().range([0, w]).padding(0.2).round(true),
        y  = d3.scaleLinear().range([h, 0]),
        yAxis = d3.axisLeft().scale(y).tickSize(-w + padl + padr),
        xAxis = d3.axisBottom().scale(x),
        visitorData = this.state.visitorData,
        pageviewData = this.state.pageviewData,
        xTick = Math.round(pageviewData.length / 7);

    var pageviewTip = d3.tip()
        .attr('class', 'd3-tip')
        .html(function(d) { return '<span>' + d.Count + '</span>' + ' pageviews' })
        .offset([-12, 0]);

    var visitorTip = d3.tip()
        .attr('class', 'd3-tip')
        .html(function(d) { return '<span>' + d.Count + '</span>' + ' visitors' })
        .offset([-12, 0]);

    var graph = d3.select('#graph');

    // remove previous graph
    graph.selectAll('*').remove();

    var vis = graph
        .append('svg')
        .attr('width', w + padl + padr)
        .attr('height', h + padt + padb)
        .append('g')
        .attr('transform', 'translate(' + padl + ',' + padt + ')');

    vis.call(pageviewTip);
    vis.call(visitorTip);

    var max = d3.max(pageviewData, function(d) { return d.Count });
    x.domain(pageviewData.map((d) => d.Label))
    y.domain([0, (max * 1.1)])

    var barWidth = x.bandwidth();

    // axes
    vis.selectAll('g.axis').remove();
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
    vis.selectAll('g.primary-bar').remove();
    var bars = vis.selectAll('g.primary-bar')
      .data(pageviewData)
      .enter().append('g')
      .attr('class', 'primary-bar')
      .attr('transform', function (d, i) { return "translate(" + x(d.Label) + ", 0)" });

    bars.append('rect')
      .attr('width', barWidth)
      .attr('height', (d) => (h - y(d.Count)) )
      .attr('y', (d) => y(d.Count))
      .on('mouseover', pageviewTip.show)
      .on('mouseout', pageviewTip.hide);
    
    vis.selectAll('g.sub-bar').remove();
    var visitorBars = vis.selectAll('g.sub-bar')
      .data(visitorData)
      .enter().append('g')
      .attr('class', 'sub-bar')
      .attr('transform', (d, i) => "translate(" + ( x(d.Label) + ( barWidth / 4 ) ) + ", 0)");

    visitorBars.append('rect')
      .attr('width', barWidth / 2 )
      .attr('height', (d) => (h - y(d.Count)) )
      .attr('y', (d) => y(d.Count))
      .on('mouseover', visitorTip.show)
      .on('mouseout', visitorTip.hide);

  }

  fetchData(period) {
    const before = Math.round((+new Date() ) / 1000);
    const after = before - ( period * dayInSeconds );

    // fetch visitor data
    fetch(`/api/visits/count/day?before=${before}&after=${after}`, {
      credentials: 'include'
    }).then((r) => {
      if( r.ok ) {
        return r.json();
      }
      throw new Error();
    }).then((data) => {
      this.setState({ visitorData: data })
      window.requestAnimationFrame(this.refreshChart);
    });

    // fetch pageview data
    fetch(`/api/pageviews/count/day?before=${before}&after=${after}`, {
      credentials: 'include'
    }).then((r) => {
      if( r.ok ) {
        return r.json();
      }

      throw new Error();
    }).then((data) => {
        this.setState({ pageviewData: data })
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

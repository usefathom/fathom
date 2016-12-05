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
    var padt = 20, padb = 40, padr = 40, padl = 40,
        h = 300,
        w = document.getElementById('graph').parentNode.clientWidth - padl - padr,
        x  = d3.scaleBand().range([0, w]).padding(0.2),
        y  = d3.scaleLinear().range([h, 0]),
        yAxis = d3.axisLeft().scale(y).tickSize(-w + padl + padr),
        xAxis = d3.axisBottom().scale(x),
        xTick = Math.round(this.state.visitorData.length / 7);

    var tip = d3.tip()
        .attr('class', 'd3-tip')
        .html(function(d) { return '<span>' + d.Count + '</span>' + ' visitors' })
        .offset([-12, 0]);

    var vis = d3.select('#graph')
        .append('svg')
        .attr('width', w)
        .attr('height', h + padt + padb)
        .append('g')
        .attr('transform', 'translate(' + padl + ',' + padt + ')');

    vis.call(tip);

    var max = d3.max(this.state.visitorData, function(d) { return d.Count });
    x.domain(this.state.visitorData.map((d) => d.Label))
    y.domain([0, max])

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
    var data = this.state.visitorData;
    var bars = vis.selectAll('g.bar')
      .data(data)
      .enter().append('g')
      .attr('class', 'bar')
      .attr('transform', function (d, i) { return "translate(" + x(d.Label) + ", 0)" });

    bars.append('rect')
      .attr('width', () => x.bandwidth() )
      .attr('height', (d) => (h - y(d.Count)) )
      .attr('y', (d) => y(d.Count))
      .on('mouseover', tip.show)
      .on('mouseout', tip.hide)
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
      window.requestAnimationFrame(this.refreshChart.bind(this));
    });

    // // fetch pageview data
    // fetch(`/api/pageviews/count/day?before=${before}&after=${after}`, {
    //   credentials: 'include'
    // }).then((r) => {
    //   if( r.ok ) {
    //     return r.json();
    //   }
    //
    //   throw new Error();
    // }).then((data) => {
    //     this.setState({ pageviewData: data })
    //     window.requestAnimationFrame(this.refreshChart.bind(this));
    // });
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

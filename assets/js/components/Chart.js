'use strict';

import { h, Component } from 'preact';
import Client from '../lib/client.js';
import { bind } from 'decko';

import * as d3 from 'd3';
import 'd3-transition';


function padZero(s) {
  return s < 10 ? "0" + s : s;
}

function padData(startUnix, endUnix, data) {
  let startDate = new Date(startUnix * 1000);
  let endDate = new Date(endUnix * 1000);
  let datamap = [];
  let newData = [];

  data.forEach((d) => {
    d.Date = d.Date.substring(0, 10);
    datamap[d.Date] = d;
  });

  while(startDate < endDate) {
    let date = startDate.getFullYear() + "-" + padZero(startDate.getMonth() + 1) + "-" + padZero(startDate.getDate());
    let data = datamap[date] ? datamap[date] : {
        "Date": date,
        "Pageviews": 0,
        "Visitors": 0,
     };

    newData.push(data);  
    startDate.setDate(startDate.getDate() + 1);
  }

 return newData;
}



class Chart extends Component {
  constructor(props) {
    super(props)

    this.state = {
      loading: false,
      data: [],
    }
  }

  componentWillReceiveProps(newProps, prevState) {
    if(newProps.before == prevState.before && newProps.after == prevState.after) {
      return;
    }

    this.fetchData(newProps.before, newProps.after);
  }


  @bind
  redrawChart() {
     let padding = { top: 10, right: 24, bottom: 48, left: 40 },
      h = 240,
      w = this.base.clientWidth;

    let height = h, width = w;
    let innerWidth = width - padding.left - padding.right;
    let innerHeight = height - padding.top - padding.bottom;

    if( this.previousGraph ) {
      this.previousGraph.selectAll('*').remove();
    } 

    let graph = this.previousGraph = d3.select(this.base)
    graph
      .append('svg').attr('width', width)
      .attr('height', height)
      .append('g').attr('transform', 'translate(' + padding.left + ', 0)');
    graph = graph.select('g');

    var t = d3.transition().duration(500).ease(d3.easeQuadOut);
    var data = this.state.data;
    var max = d3.max(data, (d) => d.Pageviews);
   
    var x  = d3.scaleBand().range([0, innerWidth]).padding(0.1),
      y  = d3.scaleLinear().range([innerHeight, 0]),
      yAxis = d3.axisLeft().scale(y).ticks(3).tickSize(-innerWidth),
      xAxis = d3.axisBottom().scale(x);

    x.domain(data.map((d) => d.Date))
    y.domain([0, (max*1.155)])

    // axes
    graph.append("g")
      .attr("class", "y axis")
      .call(yAxis);

    let nxTicks = Math.max(1, Math.round(data.length / 60));  
    let nxLabels = Math.max(1, Math.round(data.length / 15));
    let xTicks = graph.append("g")
      .attr("class", "x axis")
      .attr('transform', 'translate(0,' + innerHeight + ')')
      .call(xAxis)
    xTicks.selectAll('g text').style('display', (d, i) => { 
      return i % nxLabels != 0 ? 'none' : 'block'
    }).attr("transform", "rotate(-50)").style("text-anchor", "end");
    xTicks.selectAll('g').style('display', (d, i) => {
      return i % nxTicks != 0 ? 'none' : 'block';
    });
 
    // pageview bars
    var pageviewBars = graph.selectAll('g.pageviews')
      .data(data)
      .enter()
      .append('g')
      .attr('class', 'pageviews') 
      .attr('transform', function (d, i) { return "translate(" + x(d.Date) + ", 0)" });
      
    pageviewBars.append('rect')
      .attr('width', x.bandwidth() * 0.5)
      .attr("y", innerHeight)
      .attr("height", 0)
      .transition(t)
      .attr('y', d => y(d.Pageviews))
      .attr('height', (d) => innerHeight - y(d.Pageviews) )
      

    // visitors  
    var visitorBars = graph.selectAll('g.visitors')
      .data(data)
      .enter()
      .append('g')
      .attr('class', 'visitors')
      .attr('transform', function (d, i) { return "translate(" + ( x(d.Date) + 0.5 * x.bandwidth() ) + ", 0)" });
    
    visitorBars.append('rect')
      .attr('width', x.bandwidth() * 0.5)
      .attr("y", innerHeight)
      .attr("height", 0)
      .transition(t)
      .attr('height', (d) => (innerHeight - y(d.Visitors)) )
      .attr('y', (d) => y(d.Visitors))
  
  }

  @bind
  fetchData(before, after) {
    this.setState({ loading: true })

    Client.request(`/stats/site/groupby/day?before=${before}&after=${after}`)
      .then((d) => { 
        // request finished; check if timestamp range is still the one user wants to see
        if( this.props.before != before || this.props.after != after ) {
          return;
        }

        this.setState({ 
          loading: false,
          data: padData(after, before, d),
        });

        this.redrawChart();
      })
  }
 
  render(props, state) {
    return (
       <div id="chart"></div>
    )
  }
}

export default Chart

'use strict';

import { h, Component } from 'preact';
import Client from '../lib/client.js';
import { bind } from 'decko';

import * as d3 from 'd3';
import 'd3-transition';
d3.tip = require('d3-tip');

const formatDay = d3.timeFormat("%e"),
    formatMonth = d3.timeFormat("%b"),
    formatMonthDay = d3.timeFormat("%b %e"),
    formatYear = d3.timeFormat("%Y");

const t = d3.transition().duration(500).ease(d3.easeQuadOut);

// tooltip
const tip = d3.tip().attr('class', 'd3-tip').html((d) => (`
  <div class="tip-heading">${d.Date.toLocaleDateString()}</div>
  <div class="tip-content">
    <div class="tip-pageviews">
      <div class="tip-number">${d.Pageviews}</div>
      <div class="tip-metric">Pageviews</div>
    </div>
    <div class="tip-visitors">
      <div class="tip-number">${d.Visitors}</div>
      <div class="tip-metric">Visitors</div>
    </div>
  </div>
`));

function padZero(s) {
  return s < 10 ? "0" + s : s;
}

function timeFormatPicker(n) {
  return function(d, i) {
    if(d.getDate() === 1) {
      return d.getMonth() === 0 ? formatYear(d) : formatMonth(d) 
    } 

    if(i === 0) {
      return formatMonthDay(d)
    } else if(n < 32) {
      return formatDay(d);
    }

    return '';
  }
}

function prepareData(startUnix, endUnix, data) {
  // add timezone offset back in to get local start date
  const timezoneOffset = (new Date()).getTimezoneOffset() * 60;
  let startDate = new Date((startUnix + timezoneOffset) * 1000);
  let endDate = new Date((endUnix+timezoneOffset) * 1000);
  let datamap = [];
  let newData = [];

   // create keyed array for quick date access
  let length = data.length;
  let d, dateParts, date, key;
  for(var i=0;i<length;i++) {
    d = data[i];
    // replace date with actual date object & store in datamap
    dateParts = d.Date.split('T')[0].split('-');
    date = new Date(dateParts[0], dateParts[1]-1, dateParts[2], 0, 0, 0)
    key = date.getFullYear() + "-" + padZero(date.getMonth() + 1) + "-" + padZero(date.getDate());
    d.Date = date;
    datamap[key] = d;
  }

  // make sure we have values for each date
  let currentDate = startDate;
  while(currentDate < endDate) {
    key = currentDate.getFullYear() + "-" + padZero(currentDate.getMonth() + 1) + "-" + padZero(currentDate.getDate());
    data = datamap[key] ? datamap[key] : {
        "Pageviews": 0,
        "Visitors": 0,
        "Date": new Date(currentDate),
     };

    newData.push(data);  
    currentDate.setDate(currentDate.getDate() + 1);
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

  componentWillReceiveProps(newProps) {
    if(newProps.before == this.props.before && newProps.after == this.props.after) {
      return;
    }

    this.fetchData(newProps.before, newProps.after);
  }

  @bind
  prepareChart() {
    let padding = { top: 12, right: 12, bottom: 24, left: 40 };
    let height = 240;
    let width = this.base.clientWidth;

    this.innerWidth = width - padding.left - padding.right;
    this.innerHeight = height - padding.top - padding.bottom;

    this.ctx =  d3.select(this.base)
      .append('svg')
      .attr('width', width)
      .attr('height', height)
      .append('g')
      .attr('transform', 'translate(' + padding.left + ', '+padding.top+')')

    this.x = d3.scaleBand().range([0, this.innerWidth]).padding(0.1)
    this.y = d3.scaleLinear().range([this.innerHeight, 0])
    this.ctx.call(tip)
  }

  @bind
  redrawChart() {
    const data = this.state.data;
    this.base.parentNode.style.display = data.length <= 1 ? 'none' : '';
    if(data.length <= 1) {
      return;
    }

    if( ! this.ctx ) {
      this.prepareChart()
    }

    let graph = this.ctx;
    let innerWidth = this.innerWidth
    let innerHeight = this.innerHeight
    const max = d3.max(data, (d) => d.Pageviews);
    let x = this.x.domain(data.map((d) => d.Date))
    let y = this.y.domain([0, (max*1.1)])
    let yAxis = d3.axisLeft().scale(y).ticks(3).tickSize(-innerWidth)
    let xAxis = d3.axisBottom().scale(x).tickFormat(timeFormatPicker(data.length))

    window.xAxis = xAxis;
    window.data = data;
    
    // empty previous graph
    graph.selectAll('*').remove()

    // add axes
    let yTicks = graph.append("g")
      .attr("class", "y axis")
      .call(yAxis);
    let xTicks = graph.append("g")
      .attr("class", "x axis")
      .attr('transform', 'translate(0,' + innerHeight + ')')
      .call(xAxis)

    // hide all "day" ticks if we're watching more than 100 days of data
    if(data.length > 100) {
      xTicks.selectAll('g').filter(d => d.getDate() > 1).remove()
    }

    // add data for each day
    let days = graph.selectAll('g.day').data(data).enter()
      .append('g')
      .attr('class', 'day') 
      .attr('transform', function (d, i) { return "translate(" + x(d.Date) + ", 0)" })
      .on('mouseover', tip.show)
      .on('mouseout', tip.hide)

    let pageviews = days.append('rect')
      .attr('class', 'bar-pageviews') 
      .attr('width', x.bandwidth() * 0.5)
      .attr("y", innerHeight)
      .attr("height", 0)
      .transition(t)
      .attr('y', d => y(d.Pageviews))
      .attr('height', (d) => innerHeight - y(d.Pageviews))

    let visitors = days.append('rect')
      .attr('class', 'bar-visitors')
      .attr('width', x.bandwidth() * 0.5)
      .attr("y", innerHeight)
      .attr("height", 0)
      .attr('transform', function (d, i) { return "translate(" + ( 0.5 * x.bandwidth() ) + ", 0)" })
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
          data: prepareData(after, before, d),
        })
        this.redrawChart()
      })
  }
 
  render(props, state) {
    return (
       <div id="chart" class={state.loading ? 'loading': ''}></div>
    )
  }
}

export default Chart

'use strict';

import { h, Component } from 'preact';
import Client from '../lib/client.js';
import { bind } from 'decko';
import * as numbers from '../lib/numbers.js';
import * as d3 from 'd3';
import 'd3-transition';
d3.tip = require('d3-tip');

const formatMonth = d3.timeFormat("%b"),
  formatMonthDay = d3.timeFormat("%b %e");

const t = d3.transition().duration(600).ease(d3.easeQuadOut);
const xTickFormat = (len) => {
  return {
    hour: (d, i) => {
      if(len <= 24 && d.getHours() == 0 || d.getHours() == 12) {
        return d.getHours() + ":00";
      }

      if(i === 0 || i === len-1) {
        return formatMonthDay(d);
      }

      return '';
    },

    day: (d, i) => {
      if(i === 0 || i === len-1) {
        return formatMonthDay(d);
      }

      return '';
    },
    month: (d, i) => {
      if(len>24) {
        return d.getFullYear();
      }

      return d.getMonth() === 0 ? d.getFullYear() : formatMonth(d);
    }
  }
}

class Chart extends Component {
  constructor(props) {
    super(props)

    this.state = {
      loading: false,
      data: [],
      chartData: [],
      diffInDays: 1,
    }
  }

  componentWillReceiveProps(newProps) {
    let daysDiff = Math.round((newProps.dateRange[1]-newProps.dateRange[0])/1000/24/60/60);

    this.setState({
      diffInDays: daysDiff,
      tickStep: newProps.tickStep,
    })

    if( newProps.siteId != this.props.siteId || newProps.dateRange[0] != this.props.dateRange[0] || newProps.dateRange[1] != this.props.dateRange[1] ) {
      this.fetchData(newProps)
    } else if (newProps.tickStep != this.props.tickStep) {
      this.chartData()
      this.redrawChart()
    }
  }

  @bind
  chartData() {
    let startDate = new Date(this.props.dateRange[0]);
    let endDate = this.props.dateRange[1];
    let newData = [];

    // if grouping by month, fix date to 1st of month
    if(this.state.tickStep === 'month') {
      startDate.setDate(1);
    }

    // instantiate JS Date objects
    let data = this.state.data.map(d => {
      d.Date = new Date(d.Date);
      return d
    })
  
    // make sure we have values for each date (so 0 value for gaps)
    let currentDate = startDate, nextDate, tick, offset = 0;
    while(currentDate < endDate) {
      tick = {
          "Pageviews": 0,
          "Visitors": 0,
          "Date": new Date(currentDate),
      };

      nextDate = new Date(currentDate)

      switch(this.state.tickStep) {
        case 'hour':
        nextDate.setHours(nextDate.getHours() + 1);
        break;

        case 'day':
        nextDate.setDate(nextDate.getDate() + 1);
        break;

        case 'month':
        nextDate.setMonth(nextDate.getMonth() + 1);
        break;
      }

      // grab data that falls between currentDate & nextDate
      for(let i=data.length-offset-1; i>=0; i--) {
        // Because 9AM should be included in 9AM-10AM range, check for equality here
        if( data[i].Date >= nextDate) {
          break;
        }

         // increment offset so subsequent dates can skip first X items in array
         offset += 1;

        // continue to next item in array if we're still below our target date
        if( data[i].Date < currentDate) {
          continue;
        }

        // add to tick data
        tick.Pageviews += data[i].Pageviews;
        tick.Visitors += data[i].Visitors;
      }

      newData.push(tick);  
      currentDate = nextDate;
    }

    this.setState({
      chartData: newData,
    })
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

      // tooltip
    this.tip = d3.tip().attr('class', 'd3-tip').html((d) => {
      let title =  d.Date.toLocaleDateString();

      if(this.state.tickStep === 'hour') {
        title += ` ${d.Date.getHours()}:00 - ${d.Date.getHours() + 1}:00`
      } 

      return (`
      <div class="tip-heading">${title}</div>
      <div class="tip-content">
        <div class="tip-pageviews">
          <div class="tip-number">${d.Pageviews}</div>
          <div class="tip-metric">Pageviews</div>
        </div>
        <div class="tip-visitors">
          <div class="tip-number">${d.Visitors}</div>
          <div class="tip-metric">Visitors</div>
        </div>
      </div>`
      )});

    this.ctx.call(this.tip)
  }

  @bind
  redrawChart() {
    let data = this.state.chartData;

    if( ! this.ctx ) {
      this.prepareChart()
    }

    let graph = this.ctx;
    let innerWidth = this.innerWidth
    let innerHeight = this.innerHeight
    const max = d3.max(data, d => d.Pageviews); 
    let x = this.x.domain(data.map(d => d.Date))
    let y = this.y.domain([0, max*1.1])
    let yAxis = d3.axisLeft().scale(y).ticks(3).tickSize(-innerWidth).tickFormat(v => numbers.formatPretty(v))
    let xAxis = d3.axisBottom().scale(x).tickFormat(xTickFormat(data.length)[this.state.tickStep])

    // only show first and last tick if we have more than 24 ticks to show
    if(data.length > 24) {
      xAxis.tickValues(data.map(d => d.Date).filter((d, i) => i === 0 || i === data.length-1))
    }

    // empty previous graph
    graph.selectAll('*').remove()

    // add text indicating there's no data yet
    if( max === 0 ) {
      graph.append('text')
        .attr('class', 'muted')
        .attr("text-anchor", "middle")
        .attr('x', innerWidth / 2 - 30)
        .attr('y', innerHeight / 2)
        .text('Nothing here, yet.')
    }

    // add axes
    let yTicks = graph.append("g")
      .attr("class", "y axis")
      .call(yAxis);

    let xTicks = graph.append("g")
      .attr("class", "x axis")
      .attr('transform', 'translate(0,' + innerHeight + ')')
      .call(xAxis)
    
    // add data for each tick that we have something to show for
    let barWidth = x.bandwidth()
    let ticks = graph.selectAll('.item')
      .data(data.filter(d => d.Pageviews > 0 || d.Visitors > 0)).enter()
      .append('g')
      .attr('class', 'item') 
      
    let pageviews = ticks.append('rect')
      .attr('class', 'bar-pageviews') 
      .attr('x', d => x(d.Date))
      .attr('width', barWidth)
      .attr("y", innerHeight)
      .attr("height", 0)
      
    let visitors = ticks.append('rect')
      .attr('class', 'bar-visitors')
      .attr('x', d => x(d.Date) )
      .attr('width', barWidth)
      .attr("y", innerHeight)
      .attr("height", 0)
    
    pageviews.transition(t)
      .attr('y', d => y(d.Pageviews))
      .attr('height', (d) => innerHeight - y(d.Pageviews)) 

    visitors.transition(t)
      .attr('height', (d) => (innerHeight - y(d.Visitors)) )
      .attr('y', (d) => y(d.Visitors))   
      
    // add event listeners for tooltips
    ticks.on('mouseover', this.tip.show).on('mouseout', this.tip.hide)   
  }

  @bind
  fetchData(props) {
    this.setState({ loading: true })

    let before = props.dateRange[1]/1000;
    let after = props.dateRange[0]/1000;

    Client.request(`/sites/${props.siteId}/stats/site?before=${before}&after=${after}`)
      .then(data => { 

        this.setState({ 
          loading: false,
          data: data,
        })

        this.chartData()
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

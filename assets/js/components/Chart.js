'use strict';

import { h, Component } from 'preact';
import Client from '../lib/client.js';
import { bind } from 'decko';
import * as d3 from 'd3';

function padData(data) {
  for(let i=0; i<data.length; i++) {
    data[i].Date = data[i].Date.substring(0, 10);
  }
  return data;
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

  componentDidMount() {
     var padding = 20,
      h = 300,
      w = this.base.clientWidth - ( padding * 3);

      this.vis = d3.select(this.base)
      .append('svg')
      .attr('width', w + (padding * 2))
      .attr('height', h + (padding * 2))
      .append('g')
      .attr('transform', 'translate(' + (2*padding) + ',' + padding + ')');
  }

  @bind
  redrawChart() {
    var data = this.state.data;
    var max = d3.max(data, (d) => d.Pageviews);

    var h = 300, 
      padding = 20,
      w = this.base.clientWidth - ( padding * 3),
      x  = d3.scaleBand().range([0, w]).padding(0.2).round(true),
      y  = d3.scaleLinear().range([h, 0]),

      yAxis = d3.axisLeft().scale(y),
      xAxis = d3.axisBottom().scale(x);

    x.domain(data.map((d) => d.Date))
    y.domain([0, (max * 1.1)])

    // clear all previous data
    this.vis.selectAll('*').remove();

    // axes
    this.vis.append("g")
      .attr("class", "y axis")
      .call(yAxis);

    this.vis.append("g")
      .attr("class", "x axis")
      .attr('transform', 'translate(0,' + h + ')')
      .call(xAxis)
      .selectAll('g')

    // bars
    var bars = this.vis.selectAll('g.pageviews')
      .data(data)
      .enter()
      .append('g')
      .attr('class', 'pageviews')
      .attr('transform', function (d, i) { return "translate(" + x(d.Date) + ", 0)" });

    bars.append('rect')
      .attr('width', x.bandwidth())
      .attr('height', (d) => (h - y(d.Pageviews)) )
      .attr('y', (d) => y(d.Pageviews))

    // visitors  
    var visitorBars = this.vis.selectAll('g.visitors')
      .data(data)
      .enter()
      .append('g')
      .attr('class', 'visitors')
      .attr('transform', function (d, i) { return "translate(" + ( x(d.Date) + 0.25 * x.bandwidth() ) + ", 0)" });
    
    visitorBars.append('rect')
      .attr('width', x.bandwidth() * 0.5)
      .attr('height', (d) => (h - y(d.Visitors)) )
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
          data: padData(d),
        });

        this.redrawChart();
      })
  }
 
  render(props, state) {
    return (
       <div></div>
    )
  }
}

export default Chart

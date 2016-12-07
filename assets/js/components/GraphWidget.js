'use strict';

import { h, render, Component } from 'preact';
import Graph from './Graph.js';

class GraphWidget extends Component {

  constructor(props) {
    super(props)
    this.state = {
      showPageviews: true,
      showVisitors: true,
    }
  }

  render() {
    return (
      <div class="block">
        <div class="pull-right">
          <label class="inline small-margin-right"><input type="checkbox" checked={this.state.showPageviews} onchange={(e) => this.setState({ showPageviews: e.target.checked })} /> Pageviews</label>
          <label class="inline"><input type="checkbox" checked={this.state.showVisitors} onchange={(e) => this.setState({ showVisitors: e.target.checked })} /> Visitors</label>
        </div>
        <Graph period={this.props.period} showPageviews={this.state.showPageviews} showVisitors={this.state.showVisitors} />
      </div>
    )
  }
}

export default GraphWidget

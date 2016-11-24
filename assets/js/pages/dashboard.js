'use strict'

import { h, render, Component } from 'preact';
import LogoutButton from '../components/LogoutButton.js';
import Pageviews from '../components/Pageviews.js';
import Realtime from '../components/Realtime.js';
import Graph from '../components/Graph.js';
import DatePicker from '../components/DatePicker.js';
import Table from '../components/Table.js';

class Dashboard extends Component {
  constructor(props) {
    super(props)

    this.state = {
      period: 7
    }
  }

  render() {
    return (<div class="container">
        <header class="header cf">
          <h1 class="pull-left">Ana <small>open web analytics</small></h1>
          <div class="pull-right">
            <LogoutButton onSuccess={this.props.onLogout} />
          </div>
        </header>
        <Realtime />
        <div class="clear">
          <DatePicker period={this.state.period} onChoose={(p) => { this.setState({ period: p })}} />
        </div>
        <Graph period={this.state.period} />
        <div>
          <Pageviews period={this.state.period} />
          <Table period={this.state.period} endpoint="languages" title="Languages" headers={["#", "Language", "Count", "%"]} />
          <Table period={this.state.period} endpoint="screen-resolutions" title="Screen Resolutions" headers={["#", "Screen Resolution", "Count", "%"]} />
          <Table period={this.state.period} endpoint="countries" title="Countries" headers={["#", "Country", "Count", "%"]} />
          <Table period={this.state.period} endpoint="browsers" title="Browsers" headers={["#", "Browser", "Count", "%"]} />
        </div>
    </div>)
  }
}

export default Dashboard

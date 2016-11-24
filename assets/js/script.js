'use strict';

import { h, render, Component } from 'preact';
import LoginForm from './components/LoginForm.js';
import LogoutButton from './components/LogoutButton.js';
import Pageviews from './components/Pageviews.js';
import Realtime from './components/Realtime.js';
import Graph from './components/Graph.js';
import DatePicker from './components/DatePicker.js';
import Table from './components/Table.js';

class App extends Component {
  constructor(props) {
    super(props)

    this.state = window.state = {
      authenticated: document.cookie.indexOf('auth') > -1,
      period: 7
    }
  }

  render() {

    // logged-in
    if( this.state.authenticated ) {
      return (
        <div class="container">
            <header class="header cf">
              <h1 class="pull-left">Ana <small>open web analytics</small></h1>
              <div class="pull-right">
                <LogoutButton onSuccess={() => { this.setState({ authenticated: false })}} />
              </div>
            </header>
            <Realtime />
            <div class="clear"><DatePicker period={this.state.period} onChoose={(p) => { this.setState({ period: p })}} /></div>
            <Graph period={this.state.period} />
            <div class="masonry-container">
              <Pageviews period={this.state.period} />
              <Table period={this.state.period} endpoint="languages" title="Languages" headers={["#", "Language", "Count", "%"]} />
              <Table period={this.state.period} endpoint="screen-resolutions" title="Screen Resolutions" headers={["#", "Screen Resolution", "Count", "%"]} />
              <Table period={this.state.period} endpoint="countries" title="Countries" headers={["#", "Country", "Count", "%"]} />
              <Table period={this.state.period} endpoint="browsers" title="Browsers" headers={["#", "Browser", "Count", "%"]} />
            </div>
        </div>
      )
    }

    // logged-out
    return (
        <div class="container">
          <header class="header cf">
            <h1 class="pull-left">Ana <small>open web analytics</small></h1>
          </header>
          <LoginForm onAuth={() => { this.setState({ authenticated: true })}} />
        </div>
    )
  }
}

render(<App />, document.getElementById('root'));

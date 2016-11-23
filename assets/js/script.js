'use strict';

import { h, render, Component } from 'preact';
import LoginForm from './components/LoginForm.js';
import LogoutButton from './components/LogoutButton.js';
import Pageviews from './components/Pageviews.js';
import Realtime from './components/Realtime.js';
import Graph from './components/Graph.js';

class App extends Component {
  constructor(props) {
    super(props)

    this.state = {
      authenticated: document.cookie.indexOf('auth') > -1,
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
            <Graph />
            <Pageviews />
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

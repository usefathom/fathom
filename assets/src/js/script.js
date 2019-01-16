'use strict';

import { h, render, Component } from 'preact'
import Login from './pages/login.js'
import Dashboard from './pages/dashboard.js'
import { bind } from 'decko';
import Client from './lib/client.js';

class App extends Component {
  constructor(props) {
    super(props)

    this.state = {
      authenticated: document.cookie.indexOf('auth') > -1,
      autherized: document.cookie.indexOf('auth') > -1,
    }

    this.fetchAuthStatus()
  }

  @bind
  fetchAuthStatus() {
    Client.request(`session`, {}, true)
      .then(({ d, status }) => {
        this.setState({
          autherized: status === 200,
          authenticated: d,
        })
      })
  }

  @bind
  logout() {
    this.setState({ 
      authenticated: false,
      autherized: false,
    })
  }

  @bind
  autherized() {
    return this.state.autherized
  }

  @bind
  authenticated() {
    return this.state.authenticated
  }

  render(props, state) {
    // logged-in
    if( state.autherized || state.authenticated ) {
      return <Dashboard onLogout={this.logout} onLogin={this.fetchAuthStatus} autherized={this.autherized} authenticated={this.authenticated} />
    }

    // logged-out
    return <Login onLogin={this.fetchAuthStatus} />
  }
}

render(<App />, document.getElementById('root'));

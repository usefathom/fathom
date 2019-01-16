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
      authorized: document.cookie.indexOf('auth') > -1,
    }

    this.fetchAuthStatus()
  }

  @bind
  fetchAuthStatus() {
    Client.request(`session`, {}, true)
      .then(({ d, status }) => {
        this.setState({
          authorized: status === 200,
          authenticated: d,
        })
      })
  }

  @bind
  logout() {
    this.setState({ 
      authenticated: false,
      authorized: false,
    })
  }

  @bind
  authorized() {
    return this.state.authorized
  }

  @bind
  authenticated() {
    return this.state.authenticated
  }

  render(props, state) {
    // logged-in
    if( state.authorized || state.authenticated ) {
      return <Dashboard onLogout={this.logout} onLogin={this.fetchAuthStatus} authorized={this.authorized} authenticated={this.authenticated} />
    }

    // logged-out
    return <Login onLogin={this.fetchAuthStatus} />
  }
}

render(<App />, document.getElementById('root'));

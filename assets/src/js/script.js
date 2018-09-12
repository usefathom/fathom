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
      authenticated: document.cookie.indexOf('auth') > -1
    }

    this.fetchAuthStatus()
  }

  @bind
  fetchAuthStatus() {
    Client.request(`session`)
      .then((d) => { 
        this.setState({ authenticated: d })
      })
  }

  @bind
  toggleAuth() {
    this.setState({ 
      authenticated: !this.state.authenticated 
    })
  }

  render(props, state) {
    // logged-in
    if( state.authenticated ) {
      return <Dashboard onLogout={this.toggleAuth} />
    }

    // logged-out
    return <Login onLogin={this.toggleAuth} />
  }
}

render(<App />, document.getElementById('root'));

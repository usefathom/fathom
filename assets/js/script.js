'use strict';

import { h, render, Component } from 'preact'
import Login from './pages/login.js'
import Dashboard from './pages/dashboard.js'
import { bind } from 'decko';

class App extends Component {
  constructor(props) {
    super(props)

    this.state = {
      authenticated: document.cookie.indexOf('auth') > -1
    }
  }

  @bind
  toggleAuth() {
    this.setState({ 
      authenticated: !this.state.authenticated 
    })
  }

  render() {
    // logged-in
    if( this.state.authenticated ) {
      return <Dashboard onLogout={this.toggleAuth} />
    }

    // logged-out
    return <Login onLogin={this.toggleAuth} />
  }
}

render(<App />, document.getElementById('root'));

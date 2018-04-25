'use strict';

import { h, render, Component } from 'preact'
import Login from './pages/login.js'
import Dashboard from './pages/dashboard.js'

class App extends Component {
  constructor(props) {
    super(props)

    this.state = {
      authenticated: document.cookie.indexOf('auth') > -1
    }
  }

  render() {
    // logged-in
    if( this.state.authenticated ) {
      return <Dashboard onLogout={() => this.setState({ authenticated: false })} />
    }

    // logged-out
    return <Login onLogin={() => this.setState({ authenticated: true })} />
  }
}

render(<App />, document.getElementById('root'));

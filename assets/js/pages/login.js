'use strict';

import { h, render, Component } from 'preact';
import LoginForm from '../components/LoginForm.js';

class Login extends Component {
  render() {
    return (
      <div class="container">
        <header class="header cf">
          <h1 class="pull-left">Ana <small>open web analytics</small></h1>
        </header>
        <LoginForm onSuccess={this.props.onLogin} />
      </div>
    )
  }
}

export default Login

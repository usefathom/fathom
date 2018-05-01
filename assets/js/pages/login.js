'use strict';

import { h, render, Component } from 'preact';
import LoginForm from '../components/LoginForm.js';

class Login extends Component {
  render() {
    return (
      <div>
        <div class="container">
          <LoginForm onSuccess={this.props.onLogin}/>
        </div>
      </div>
    )
  }
}

export default Login

'use strict';

import { h, render, Component } from 'preact';
import LoginForm from '../components/LoginForm.js';

class Login extends Component {
  render(props, state) {
    return (
      <div class="flex-rapper login-page">
        <div class="login-rapper">
          <LoginForm onSuccess={props.onLogin} />
          <small><a href="https://usefathom.com">Fathom Analytics</a>{/* &middot; <a href="#lost">Password reset</a> */}</small>
        </div>
      </div>
    )
  }
}

export default Login

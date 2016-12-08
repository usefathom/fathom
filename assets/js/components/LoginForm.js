'use strict';

import { h, render, Component } from 'preact';
import Client from '../lib/client.js';

class LoginForm extends Component {

  constructor(props) {
    super(props)
    this.handleSubmit = this.handleSubmit.bind(this);
    this.state = {
      email: '',
      password: '',
    }
  }

  handleSubmit(e) {
    e.preventDefault();

    Client.request('session', {
      method: "POST",
      data: {
        email: this.state.email,
        password: this.state.password,
      }
    }).then((r) => { this.props.onSuccess() })
  }

  render() {
    return (
      <div class="block">
        <h2>Login</h2>
        <p>Please enter your credentials to access your Ana dashboard.</p>
        <form method="POST" onSubmit={this.handleSubmit}>
          <div class="small-margin">
            <label>Email address</label>
            <input type="email" name="email" placeholder="Email address" onChange={this.linkState('email')} required="required" />
          </div>
          <div class="small-margin">
            <label>Password</label>
            <input type="password" name="password" placeholder="**********" onChange={this.linkState('password')} required="required" />
          </div>
          <div class="small-margin">
            <input type="submit" value="Sign in" />
          </div>
        </form>
      </div>
    )
  }
}

export default LoginForm

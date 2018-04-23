'use strict';

import { h, render, Component } from 'preact';
import Client from '../lib/client.js';
import Notification from '../components/Notification.js';
import { bind, memoize, debounce } from 'decko';

class LoginForm extends Component {

  constructor(props) {
    super(props)
    this.handleSubmit = this.handleSubmit.bind(this);
    this.state = {
      email: '',
      password: '',
      message: ''
    }
  }

  @bind
  handleSubmit(e) {
    e.preventDefault();

    Client.request('session', {
      method: "POST",
      data: {
        email: this.state.email,
        password: this.state.password,
      }
    }).then((r) => {
        this.props.onSuccess()
    }).catch((e) => {
      this.setState({
        message: e.message,
        password: ''
      })

    })
  }

  @bind
  updatePassword(e) {
    this.setState({ password: e.target.value });
  }

  @bind
  updateEmail(e) {
    this.setState({ email: e.target.value });
  }

  render() {
    return (
      <div class="block">
        <h2>Login</h2>
        <p>Please enter your login credentials to access your Ana dashboard.</p>
        <form method="POST" onSubmit={this.handleSubmit}>
          <div class="small-margin">
            <label>Email address</label>
            <input type="email" name="email" placeholder="Email address" value={this.state.email} onInput={this.updateEmail} required="required" />
          </div>
          <div class="small-margin">
            <label>Password</label>
            <input type="password" name="password" placeholder="**********" value={this.state.password} onInput={this.updatePassword} required="required" />
          </div>
          <div class="small-margin">
            <input type="submit" value="Sign in" />
          </div>
        </form>
        <Notification message={this.state.message} kind="" />
      </div>

    )
  }
}

export default LoginForm

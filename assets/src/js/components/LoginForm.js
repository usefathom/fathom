'use strict';

import { h, render, Component } from 'preact';
import Client from '../lib/client.js';
import Notification from '../components/Notification.js';
import { bind } from 'decko';

class LoginForm extends Component {

  constructor(props) {
    super(props)
    this.state = {
      email: '',
      password: '',
      message: ''
    }
  }

  @bind
  handleSubmit(e) {
    e.preventDefault();
    this.setState({ message: '' });

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
        message: e.code === 'invalid_credentials' ? "Invalid username or password" : e.message,
        password: ''
      });
    });
  }

  @bind
  updatePassword(e) {
    this.setState({ password: e.target.value });
  }

  @bind
  updateEmail(e) {
    this.setState({ email: e.target.value });
  }

  @bind
  clearMessage() {
    this.setState({ message: '' });
  }

  render(props, state) {
    return (
      <form method="POST" onSubmit={this.handleSubmit}>
        <div class="">
          <label><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 496 512"><path d="M248 8C111 8 0 119 0 256s111 248 248 248 248-111 248-248S385 8 248 8zm128 421.6c-35.9 26.5-80.1 42.4-128 42.4s-92.1-15.9-128-42.4V416c0-35.3 28.7-64 64-64 11.1 0 27.5 11.4 64 11.4 36.6 0 52.8-11.4 64-11.4 35.3 0 64 28.7 64 64v13.6zm30.6-27.5c-6.8-46.4-46.3-82.1-94.6-82.1-20.5 0-30.4 11.4-64 11.4S204.6 320 184 320c-48.3 0-87.8 35.7-94.6 82.1C53.9 363.6 32 312.4 32 256c0-119.1 96.9-216 216-216s216 96.9 216 216c0 56.4-21.9 107.6-57.4 146.1zM248 120c-48.6 0-88 39.4-88 88s39.4 88 88 88 88-39.4 88-88-39.4-88-88-88zm0 144c-30.9 0-56-25.1-56-56s25.1-56 56-56 56 25.1 56 56-25.1 56-56 56z"/></svg></label>
          <input type="email" name="email" placeholder="Email address" required="" value={state.email} onInput={this.updateEmail}  />
         </div>
        
        <div class="">
          <label><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512"><path d="M224 420c-11 0-20-9-20-20v-64c0-11 9-20 20-20s20 9 20 20v64c0 11-9 20-20 20zm224-148v192c0 26.5-21.5 48-48 48H48c-26.5 0-48-21.5-48-48V272c0-26.5 21.5-48 48-48h16v-64C64 71.6 136-.3 224.5 0 312.9.3 384 73.1 384 161.5V224h16c26.5 0 48 21.5 48 48zM96 224h256v-64c0-70.6-57.4-128-128-128S96 89.4 96 160v64zm320 240V272c0-8.8-7.2-16-16-16H48c-8.8 0-16 7.2-16 16v192c0 8.8 7.2 16 16 16h352c8.8 0 16-7.2 16-16z"/></svg></label>
          <input type="password" name="password" placeholder="**********" required="" autocomplete="off" value={state.password} onInput={this.updatePassword} />
        </div>
         
        <div><button type="submit">Sign in</button></div>

        <Notification message={state.message} kind="" onDismiss={this.clearMessage} />
      </form>
    )
  }
}

export default LoginForm

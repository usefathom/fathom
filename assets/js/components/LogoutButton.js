'use strict';

import { h, render, Component } from 'preact';
import Client from '../lib/client.js';
import { bind } from 'decko';

class LogoutButton extends Component {

  @bind
  handleSubmit(e) {
    e.preventDefault();

    Client.request('session', {
      method: "DELETE",
    }).then((r) => { this.props.onSuccess() })
  }

  render() {
    return (
      <a href="#" onClick={this.handleSubmit}>Sign out</a>
    )
  }
}

export default LogoutButton

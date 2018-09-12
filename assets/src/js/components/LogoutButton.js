'use strict';

import { h, Component } from 'preact';
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
    if(document.cookie.indexOf('auth') < 0) {
      return ''
    }

    return (
      <a href="#" onClick={this.handleSubmit}>Sign out</a>
    )
  }
}

export default LogoutButton

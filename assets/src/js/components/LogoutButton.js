'use strict';

import { h, render, Component } from 'preact';
import Client from '../lib/client.js';

class LogoutButton extends Component {

  constructor(props) {
    super(props)
    this.handleSubmit = this.handleSubmit.bind(this);
  }

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

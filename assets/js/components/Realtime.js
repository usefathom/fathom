'use strict';

import { h, render, Component } from 'preact';
import Client from '../lib/client.js';
import { bind } from 'decko';

class Realtime extends Component {

  constructor(props) {
    super(props)

    this.state = {
      count: 0
    }
    this.fetchData();
    window.setInterval(this.fetchData, 15000);
  }

  @bind
  fetchData() {
    Client.request(`visitors/count/realtime`)
      .then((d) => { this.setState({ count: d })})
  }

  render() {
    let visitors = this.state.count == 1 ? 'visitor' : 'visitors';
    return (
        <span><span class="count">{this.state.count}</span> <span>current {visitors}</span></span>
    )
  }
}

export default Realtime

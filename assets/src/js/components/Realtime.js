'use strict';

import { h, render, Component } from 'preact';
import Client from '../lib/client.js';

class Realtime extends Component {

  constructor(props) {
    super(props)

    this.state = {
      count: 0
    }
    this.fetchData = this.fetchData.bind(this);
    this.fetchData();
    window.setInterval(this.fetchData, 15000);
  }

  fetchData() {
    Client.request(`visitors/count/realtime`)
      .then((d) => { this.setState({ count: d })})
  }

  render() {
    let visitors = this.state.count == 1 ? 'visitor' : 'visitors';
    return (
      <div class="block block-float">
        <span class="count">{this.state.count}</span> <span>{visitors} on the site right now.</span>
      </div>
    )
  }
}

export default Realtime

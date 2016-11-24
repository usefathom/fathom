'use strict';

import { h, render, Component } from 'preact';

class Realtime extends Component {

  constructor(props) {
    super(props)

    this.state = {
      count: 0
    }
    this.fetchData = this.fetchData.bind(this);
    this.fetchData();
  }

  fetchData() {
    return fetch('/api/visits/count/realtime', {
      credentials: 'include'
    })
      .then((r) => r.json())
      .then((data) => {
        this.setState({ count: data })
    });
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

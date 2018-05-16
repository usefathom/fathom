'use strict';

import { h, Component } from 'preact';
import Client from '../lib/client.js';
import { bind } from 'decko';

class Realtime extends Component {

  constructor(props) {
    super(props)

    this.state = {
      count: 0
    }
  }

  componentDidMount() {
      this.fetchData();
      this.interval = window.setInterval(this.fetchData, 15000);
  }

  componentWillUnmount() {
      clearInterval(this.interval);
  }

  @bind
  fetchData() {
    Client.request(`stats/site/realtime`)
      .then((d) => { 
        this.setState({ count: d })
        document.title = ( d > 0 ? d + ' current visitors — Fathom' : 'Fathom' );
      })
      .catch((e) => {
        if(e.message == 401) {
          this.props.onError();
        }
      })
  }

  render(props, state) {
    let visitorText = state.count == 1 ? 'visitor' : 'visitors';
    return (
        <span><span class="count">{state.count}</span> <span>current {visitorText}</span></span>
    )
  }
}

export default Realtime

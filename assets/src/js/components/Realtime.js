'use strict';

import { h, Component } from 'preact';
import Client from '../lib/client.js';
import { bind } from 'decko';
import * as numbers from '../lib/numbers.js';

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
  setDocumentTitle() {
    // update document title
    let visitorText = this.state.count == 1 ? 'visitor' : 'visitors';
    document.title = ( this.state.count > 0 ? `${numbers.formatPretty(this.state.count)} current ${visitorText} — Fathom` : 'Fathom' );
  }

  @bind
  fetchData() {
    Client.request(`/sites/${this.props.site.id}/stats/site/realtime`)
      .then((d) => { 
        this.setState({ count: d })
        this.setDocumentTitle();
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
        <span><span class="count">{numbers.formatPretty(state.count)}</span> <span>current {visitorText}</span></span>
    )
  }
}

export default Realtime

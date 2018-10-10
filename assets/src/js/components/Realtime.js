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
      this.fetchData(this.props.siteId);
      this.interval = window.setInterval(this.handleIntervalEvent, 15000);
  }

  componentWillUnmount() {
      window.clearInterval(this.interval);
  }

  componentWillReceiveProps(newProps, newState) {
    if(!this.paramsChanged(this.props, newProps)) {
      return;
    }

    this.fetchData(newProps.siteId)
  }

  paramsChanged(o, n) {
    return o.siteId != n.siteId;
  }

  @bind
  setDocumentTitle() {
    // update document title
    let visitorText = this.state.count == 1 ? 'visitor' : 'visitors';
    document.title = ( this.state.count > 0 ? `${numbers.formatPretty(this.state.count)} current ${visitorText} — Fathom` : 'Fathom' );
  }

  @bind 
  handleIntervalEvent() {
    this.fetchData(this.props.siteId)
  }

  @bind
  fetchData(siteId) {
    let url = `/sites/${siteId}/stats/site/realtime`
    Client.request(url)
      .then((d) => { 
        this.setState({ count: d })
        this.setDocumentTitle();
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

'use strict';

import { h, Component } from 'preact';
import Client from '../lib/client.js';
import { bind } from 'decko';
import CountWidget from './CountWidget.js';


class Sidebar extends Component {
  constructor(props) {
    super(props)

    this.state = {
      data: {},
      loading: false,
    }
  }

  componentWillReceiveProps(newProps, newState) {
    if(!this.paramsChanged(this.props, newProps)) {
      return;
    }

    this.fetchData(newProps);
  }

  paramsChanged(o, n) {
    return o.siteId != n.siteId || o.before != n.before || o.after != n.after;
  }

  @bind
  fetchData(props) {
    this.setState({ loading: true })

    Client.request(`/sites/${props.siteId}/stats/site/agg?before=${props.before}&after=${props.after}`)
      .then((data) => { 
        // request finished; check if timestamp range is still the one user wants to see
        if(this.paramsChanged(props, this.props)) {
          return;
        }

        // Make sure we always show at least 1 visitor when there are pageviews
        if ( data.Visitors == 0 && data.Pageviews > 0 ) {
          data.Visitors = 1
        }

        this.setState({ 
          loading: false,
          data: data
        })
      })
  }

  render(props, state) {
    return (
      <div class="box box-totals">
        <CountWidget title="Unique visitors" value={state.data.Visitors} loading={state.loading} />
        <CountWidget title="Pageviews" value={state.data.Pageviews} loading={state.loading} />
        <CountWidget title="Avg time on site" value={state.data.AvgDuration} format="duration" loading={state.loading} />
        <CountWidget title="Bounce rate" value={state.data.BounceRate} format="percentage" loading={state.loading} />
      </div>
    )
  }
}

export default Sidebar

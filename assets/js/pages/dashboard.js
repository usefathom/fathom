'use strict'

import { h, render, Component } from 'preact';
import Pageviews from '../components/Pageviews.js';
import Realtime from '../components/Realtime.js';
import GraphWidget from '../components/GraphWidget.js';
import DatePicker from '../components/DatePicker.js';
import Table from '../components/Table.js';
import HeaderBar from '../components/HeaderBar.js';
import CountWidget from '../components/CountWidget.js';

function removeImgElement(e) {
  e.target.parentNode.removeChild(e.target);
}

function formatCountryLabel(p) {
  const src = "/static/img/country-flags/"+ p.Label.toLowerCase() +".png"

  return (
    <td>
      {p.Label}
      <img height="12" src={src} class="pull-right" onError={removeImgElement} />
    </td>
  )
}

class Dashboard extends Component {
  constructor(props) {
    super(props)

    this.state = {
      period: parseInt(window.location.hash.substring(2)) || 7
    }

    this.onPeriodChoose = this.onPeriodChoose.bind(this)
  }

  onPeriodChoose(p) {
    this.setState({ period: p })
    window.history.replaceState(this.state, null, `#!${p}`)
  }

  render() {
    return (
    <div>
      <HeaderBar showLogout={true} onLogout={this.props.onLogout} />
      <div class="container">
        <Realtime />
        <div class="clear">
          <DatePicker period={this.state.period} onChoose={this.onPeriodChoose} />
        </div>
        <div class="row">
          <div class="col-2">
            <CountWidget title="Visitors" endpoint="visitors" period={this.state.period} />
          </div>
          <div class="col-2">
            <CountWidget title="Pageviews" endpoint="pageviews" period={this.state.period} />
          </div>
        </div>
        <GraphWidget period={this.state.period} />
        <div class="row">
          <div class="col-4">
            <Pageviews period={this.state.period} />
          </div>
          <div class="col-2">
            <Table period={this.state.period} endpoint="languages" title="Languages" headers={["#", "Language", "Count", "%"]} />
          </div>
        </div>
        <div class="row">
          <div class="col-2">
            <Table period={this.state.period} endpoint="screen-resolutions" title="Screen Resolutions" headers={["#", "Resolution", "Count", "%"]} />
          </div>
          <div class="col-2">
            <Table period={this.state.period} endpoint="referrers" title="Referrers" headers={["#", "URL", "Count", "%"]} labelCell={(p) => ( <td><a href={p.Label}>{p.Label.substring(0, 15).replace('https://', '').replace('http://', '')}</a></td>)} />
          </div>
          <div class="col-2">
            <Table period={this.state.period} endpoint="browsers" title="Browsers" headers={["#", "Browser", "Count", "%"]} onAuthError={this.props.onLogout} />
          </div>
        </div>
      </div>
  </div>
  )}
}

export default Dashboard

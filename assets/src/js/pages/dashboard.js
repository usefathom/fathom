'use strict'

import { h, Component } from 'preact';
import LogoutButton from '../components/LogoutButton.js';
import Realtime from '../components/Realtime.js';
import DatePicker from '../components/DatePicker.js';
import Sidebar from '../components/Sidebar.js';
import Table from '../components/Table.js';
import Chart from '../components/Chart.js';
import { bind } from 'decko';

class Dashboard extends Component {
  constructor(props) {
    super(props)

    this.state = {
      period: (window.location.hash.substring(2) || 'last-7-days'),
      before: 0,
      after: 0,
    }
  }

  @bind
  changePeriod(s) {
    this.setState({ period: s.period, before: s.before, after: s.after })
    window.history.replaceState(this.state, null, `#!${s.period}`)
  }

  render(props, state) {
    return (
    <div class="app-page wrapper">

      <header class="section">
        <nav class="main-nav animated fadeInDown">
            <ul>
              <li class="logo"><a href="/">Fathom</a></li>
              <li class="visitors"><Realtime onError={props.onLogout} /></li>
              <li class="spacer">&middot;</li>
              <li class="signout"><LogoutButton onSuccess={props.onLogout} /></li>
          </ul>
        </nav>
      </header>

      <section class="section animated fadeInUp delayed_02s">
        <nav class="date-nav">
          <DatePicker onChange={this.changePeriod} value={state.period} />
        </nav>

        <div class="boxes">
          <Sidebar before={state.before} after={state.after} />

          <div class="boxes-col">
            <div class="box box-graph">
              <Chart before={state.before} after={state.after}  />
            </div>
            <div class="box box-pages animated fadeInUp delayed_04s">
              <Table endpoint="stats/pages" headers={["Top pages", "Views", "Uniques"]} before={state.before} after={state.after} />
            </div>
            <div class="box box-referrers animated fadeInUp delayed_04s">
              <Table endpoint="stats/referrers" headers={["Top referrers", "Views", "Uniques"]} before={state.before} after={state.after} showHostname="true" />
            </div>
          </div>
        </div>
      </section>

      <footer class="section"></footer>
    </div>
  )}
}

export default Dashboard

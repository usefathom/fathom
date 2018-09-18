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
      before: 0,
      after: 0,
      isPublic: document.cookie.indexOf('auth') < 0,
    }
  }

  @bind
  updateDateRange(s) {
    this.setState({ 
      before: s.before, 
      after: s.after 
    })
  }

  render(props, state) {
    // only show logout link if this dashboard is not public
    let logoutMenuItem = state.isPublic ? '' : (
      <li class="signout"><span class="spacer">&middot;</span><LogoutButton onSuccess={props.onLogout} /></li>
    );

    return (
    <div class="app-page wrapper">

      <header class="section">
        <nav class="main-nav animated fadeInDown">
            <ul>
              <li class="logo"><a href="/">Fathom</a></li>
              <li class="visitors"><Realtime /></li>
              {logoutMenuItem}
          </ul>
        </nav>
      </header>

      <section class="section animated fadeInUp delayed_02s">
        <nav class="date-nav">
          <DatePicker onChange={this.updateDateRange} />
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

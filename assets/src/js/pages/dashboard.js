'use strict'

import { h, Component } from 'preact';
import LogoutButton from '../components/LogoutButton.js';
import Realtime from '../components/Realtime.js';
import DatePicker from '../components/DatePicker.js';
import Sidebar from '../components/Sidebar.js';
import SiteSwitcher from '../components/SiteSwitcher.js';
import SiteSettings from '../components/SiteSettings.js';
import Gearwheel from '../components/Gearwheel.js';
import Table from '../components/Table.js';
import Chart from '../components/Chart.js';
import { bind } from 'decko';
import Client from '../lib/client.js';

class Dashboard extends Component {
  constructor(props) {
    super(props)

    // TODO: Fetch sites from server and populate state
    this.state = {
      before: 0,
      after: 0,
      isPublic: document.cookie.indexOf('auth') < 0,
      site: { id: 0, name: "Default site"},
      sites: [],
      settingsOpen: false,
    }
  }

  componentDidMount() {
    this.fetchSites()
  }

  @bind 
  fetchSites() {
    Client.request(`sites`)
    .then((data) => { 
      this.setState({sites: data})
    })
  }

  @bind
  changeDateRange(s) {
    this.setState({ 
      before: s.before, 
      after: s.after 
    })
  }

  @bind 
  openSiteSettings(site) {
    this.setState( { settingsOpen: true, site: site && site.hasOwnProperty('id') ? site : this.state.site })
  }

  @bind
  closeSiteSettings() {
    this.setState({settingsOpen: false})
  }

  @bind 
  changeSite(site) {
    this.setState({site: site})
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
              <SiteSwitcher sites={state.sites} selectedSite={state.site} onChange={this.changeSite} onAdd={this.openSiteSettings} />
              <Gearwheel onClick={this.openSiteSettings} visible={!state.isPublic} />
              <li class="visitors"><Realtime /></li>
          </ul>
        </nav>
      </header>

      <section class="section animated fadeInUp delayed_02s">
        <nav class="date-nav">
          <DatePicker onChange={this.changeDateRange} />
        </nav>

        <div class="boxes">
          <Sidebar site={state.site} before={state.before} after={state.after} />

          <div class="boxes-col">
            <div class="box box-graph">
              <Chart site={state.site} before={state.before} after={state.after}  />
            </div>
            <div class="box box-pages animated fadeInUp delayed_04s">
              <Table endpoint="stats/pages" headers={["Top pages", "Views", "Uniques"]} site={state.site} before={state.before} after={state.after} />
            </div>
            <div class="box box-referrers animated fadeInUp delayed_04s">
              <Table endpoint="stats/referrers" headers={["Top referrers", "Views", "Uniques"]} site={state.site} before={state.before} after={state.after} showHostname="true" />
            </div>
          </div>
        </div>

        <div class="footer hide-on-mobile">
          <p>Use <strong>ALT + arrow-key</strong> to cycle through date ranges.</p>
        </div>
      </section>

      <footer class="section"></footer>

      <SiteSettings visible={state.settingsOpen} onClose={this.closeSiteSettings} site={state.site} />
    </div>
  )}
}

export default Dashboard

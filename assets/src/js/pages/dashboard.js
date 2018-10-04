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

    // TODO: Update state after adding / deleting site
    // TODO: Update endpoints to take site ID parameter
    this.state = {
      before: 0,
      after: 0,
      isPublic: document.cookie.indexOf('auth') < 0,
      site: { id: 0 },
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
      // TODO: Account for no sites here
      // TODO: Get selected site from localstorage
      this.setState({
        sites: data, 
        site: data[0] 
      })
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
  changeSelectedSite(site) {
    this.setState({site: site})
  }

  @bind
  updateSite(site) {
    let updated = false;
    let newSites = this.state.sites.map((s) => {
      if(s.id != site.id) {
        return s;
      }

      updated = true;
      return site;
    })

    if(!updated) {
      newSites.push(site);
    }

    this.setState({sites: newSites, site: site})
  }

  @bind 
  deleteSite(site) {
    let newSites = this.state.sites.filter((s) => (s.id != site.id))
    this.setState({ sites: newSites, site: newSites[0] })
  }


  render(props, state) {
    // only show logout link if this dashboard is not public
    let logoutMenuItem = state.isPublic ? '' : (
      <li class="signout"><span class="spacer">&middot;</span><LogoutButton onSuccess={props.onLogout} /></li>
    );

    return (
  <div class="app-page ">
     <div class="wrapper animated fadeInUp delayed_02s">

      <header class="section">
        <nav class="main-nav">
            <ul>
              <li class="logo"><a href="/">Fathom</a></li>
              <SiteSwitcher sites={state.sites} selectedSite={state.site} onChange={this.changeSelectedSite} onAdd={this.openSiteSettings} showAdd={!state.isPublic}/>
              <Gearwheel onClick={this.openSiteSettings} visible={!state.isPublic} />
              <li class="visitors"><Realtime /></li>
          </ul>
        </nav>
      </header>

      <section class="section">
        <nav class="date-nav">
          <DatePicker onChange={this.changeDateRange} />
        </nav>

        <div class="boxes">
          <Sidebar site={state.site} before={state.before} after={state.after} />

          <div class="boxes-col">
            <div class="box box-graph">
              <Chart site={state.site} before={state.before} after={state.after}  />
            </div>
            <div class="box box-pages">
              <Table endpoint="stats/pages" headers={["Top pages", "Views", "Uniques"]} site={state.site} before={state.before} after={state.after} />
            </div>
            <div class="box box-referrers">
              <Table endpoint="stats/referrers" headers={["Top referrers", "Views", "Uniques"]} site={state.site} before={state.before} after={state.after} showHostname="true" />
            </div>
          </div>
        </div>

        <div class="footer hide-on-mobile">
          <p>Use <strong>ALT + arrow-key</strong> to cycle through date ranges.</p>
        </div>
      </section>

      <footer class="section"></footer>

    </div>
    <SiteSettings visible={state.settingsOpen} onClose={this.closeSiteSettings} onUpdate={this.updateSite} onDelete={this.deleteSite} site={state.site} />
  </div>
  )}
}

export default Dashboard

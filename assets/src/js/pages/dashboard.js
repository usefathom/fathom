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

const defaultSite = { 
  id: window.localStorage.getItem('site_id') || 0, 
  name: "",
};

class Dashboard extends Component {
  constructor(props) {
    super(props)

    this.state = {
      before: 0,
      after: 0,
      isPublic: document.cookie.indexOf('auth') < 0,
      site: defaultSite,
      sites: [],
      settingsOpen: false,
      addingNewSite: false,
    }
  }

  componentDidMount() {
    this.fetchSites()
  }

  @bind 
  fetchSites() {
    Client.request(`sites`)
    .then((sites) => { 
      // open site settings when there are no sites yet
      if(sites.length == 0) {
        this.showSiteSettings({ id: 0, name: "yoursite.com"})
        return;
      }

      // if there are sites, use remembered site as selected site
      let site = sites[0];
      let s = sites.find(s => s.id == defaultSite.id);
      site = s ? s : site;

      this.setState({
        sites: sites, 
        site: site,
      })
    })
  }

  @bind
  changeDateRange(s) {
    this.setState({ 
      before: s.before, 
      after: s.after,
      period: s.period,
    })
  }

  @bind 
  showSiteSettings(site) {
    site = site && site.id == 0 ? site : this.state.site;
    this.setState({ 
      settingsOpen: true, 
      site: site,
      previousSite: this.state.site,
    })
  }

  @bind
  closeSiteSettings() {
    this.setState({
      settingsOpen: false, 

      // switch back to previous site if we were showing site settings to add a new site
      site: this.state.site.id > 0 ? this.state.site : this.state.previousSite,
    })
  }

  @bind 
  changeSelectedSite(site) {
    this.setState({
      site: site,
      previousSite: this.state.site,
    })

    window.localStorage.setItem('site_id', site.id)
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
    let newSelectedSite = newSites.length > 0 ? newSites[0] : defaultSite;
    this.setState({ 
      sites: newSites, 
      site: newSelectedSite 
    })
  }

  render(props, state) {
    // only show logout link if this dashboard is not public
    let logoutMenuItem = state.isPublic ? '' : (
      <li class="signout"><span class="spacer">&middot;</span><LogoutButton onSuccess={props.onLogout} /></li>
    );

    return (
  <div class="app-page ">
     <div class={"rapper animated fadeInUp delayed_02s " + state.period }>

      <header class="section">
        <nav class="main-nav">
            <ul>
              <li class="logo"><a href="/">{state.site.name || "Fathom"}</a></li>
              <SiteSwitcher sites={state.sites} selectedSite={state.site} onChange={this.changeSelectedSite} onAdd={this.showSiteSettings} showAdd={!state.isPublic}/>
              <Gearwheel onClick={this.showSiteSettings} visible={!state.isPublic} />
              <li class="visitors"><Realtime siteId={state.site.id} /></li>
          </ul>
        </nav>
      </header>

      <section class="section">
        <nav>
          <DatePicker onChange={this.changeDateRange} />
        </nav>

        <div class="boxes">
          <Sidebar siteId={state.site.id} before={state.before} after={state.after} />

          <div class="box box-graph">
            <Chart siteId={state.site.id} before={state.before} after={state.after}  />
          </div>
          <div class="box box-pages">
            <Table endpoint="stats/pages" headers={["Top pages", "Views", "Uniques"]} siteId={state.site.id} before={state.before} after={state.after} />
          </div>
          <div class="box box-referrers">
            <Table endpoint="stats/referrers" headers={["Top referrers", "Views", "Uniques"]} siteId={state.site.id} before={state.before} after={state.after} showHostname="true" />
          </div>
        </div>

        <footer class="section">
          <div class="half">
          <nav>
            <ul>
              <li><a href="https://usefathom.com/">Fathom</a></li>
              <li><a href="https://usefathom.com/terms/">Terms of use</a></li>
              <li><a href="https://usefathom.com/privacy/">Privacy policy</a></li>
              <li><a href="https://usefathom.com/data/">Our data policy</a></li>
              <li><LogoutButton onSuccess={props.onLogout} /></li>
            </ul>
          </nav>
          <div class="hide-on-mobile">Use <strong>the arrow keys</strong> to cycle through date ranges.</div>
          </div>
        </footer>
      </section>
    </div>
    <SiteSettings visible={state.settingsOpen} onClose={this.closeSiteSettings} onUpdate={this.updateSite} onDelete={this.deleteSite} site={state.site} />
  </div>
  )}
}

export default Dashboard

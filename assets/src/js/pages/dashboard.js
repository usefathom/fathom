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
import classNames from 'classnames';

const defaultSite = { 
  id: window.localStorage.getItem('site_id') || 1, 
  name: "",
  unsaved: true,
};

class Dashboard extends Component {
  constructor(props) {
    super(props)

    this.state = {
      dateRange: [],
      groupBy: 'day',
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
        this.showSiteSettings({ id: 1, name: "yoursite.com", unsaved: true })
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
    }).catch((e) => {
      if(e.code === 'unauthorized') {
        this.props.onLogout()
      }
    })
  }

  @bind
  changeDateRange(s) {
    this.setState({ 
      dateRange: [ s.startDate, s.endDate ],
      period: s.period,
      groupBy: s.groupBy,
    })
  }

  @bind 
  showSiteSettings(site) {
    site = site && site.unsaved ? site : this.state.site;
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
      site: this.state.site.unsaved && this.state.previousSite ? this.state.previousSite : this.state.site,
    })
  }

  @bind 
  changeSelectedSite(site) {
    let newState = {
      site: site,
    }

    if(!this.state.site.unsaved) {
      newState.previousSite = this.state.site
    } 

    this.setState(newState)
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
      
      // replace site in sites array with parameter
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
     <div class={`rapper animated fadeInUp delayed_02s ${state.period} ` + classNames({ ltday: state.dateRange[1] - state.dateRange[0] < 86400 })}>

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

      <DatePicker onChange={this.changeDateRange} />

      <section class="section">
        <div class="boxes">
          <Sidebar siteId={state.site.id} dateRange={state.dateRange} />

          <div class="box box-graph">
            <Chart siteId={state.site.id} dateRange={state.dateRange} tickStep={state.groupBy} />
          </div>
          <div class="box box-pages">
            <Table endpoint="pages" headers={["Top pages", "Views", "Uniques"]} siteId={state.site.id} dateRange={state.dateRange} />
          </div>
          <div class="box box-referrers">
            <Table endpoint="referrers" headers={["Top referrers", "Views", "Uniques"]} siteId={state.site.id} dateRange={state.dateRange} showHostname="true" />
          </div>
        </div>
        <div class="notice">
          <div class="notice-text"><strong>Want more features and less maintenance?</strong> Check out the current version of Fathom Analytics and <a href="https://usefathom.com/ref/GITHUB">start your free trial today</a>.</div>
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

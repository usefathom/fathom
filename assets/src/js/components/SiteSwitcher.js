'use strict';

import { h, Component } from 'preact';
import { bind } from 'decko';

class SiteSwitcher extends Component {

  componentDidMount() {
     
  }

  componentWillUnmount() {
      
  }

  @bind 
  selectSite(evt) {
    let itemId = evt.target.getAttribute("data-id")  
    this.props.sites.some((s) => {
        if (s.id != itemId) {
            return false;
        }

        this.props.onChange(s)
        return true;
    })   
  }

  @bind 
  addSite() {
      this.props.onAdd({ id: 0, name: "New site"})
  }

  render(props, state) {
    let sites = props.sites.map((s) => (<li class="site-switch"><a href="javascript:void(0);" data-id={s.id} onClick={this.selectSite}>{s.name}</a></li>)) 
    return (
        <li class="sites">
            <a href="javascript:void(0)">{props.selectedSite.name}</a>
            <ul>
                {sites}
                <li class="add-new"><a href="javascript:void(0);" onClick={this.addSite}>+ Add another site</a></li>
            </ul>
        </li>
    )
  }
}

export default SiteSwitcher

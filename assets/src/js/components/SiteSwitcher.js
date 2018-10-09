'use strict';

import { h, Component } from 'preact';
import { bind } from 'decko';

class SiteSwitcher extends Component {
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
      this.props.onAdd({ id: 1, name: "New site", unsaved: true })
  }

  render(props, state) {
    // show nothing if there is only 1 site and no option to add additional sites
    if(!props.showAdd && props.sites.length == 1) {
        return '';
    }  

    // otherwise, render list of sites + add button
    return (
        <li class="sites">
            <a href="javascript:void(0)">{props.selectedSite.name}</a>
            <ul>
                {props.sites.map((s) => (<li class="site-switch"><a href="javascript:void(0);" data-id={s.id} onClick={this.selectSite}>{s.name}</a></li>)) }
                {props.showAdd ? (<li class="add-new"><a href="javascript:void(0);" onClick={this.addSite}>+ Add another site</a></li>) : ''}
            </ul>
        </li>
    )
  }
}

export default SiteSwitcher

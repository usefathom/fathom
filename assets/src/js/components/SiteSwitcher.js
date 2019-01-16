'use strict';

import { h, Component } from 'preact';
import { bind } from 'decko';

class SiteSwitcher extends Component {
  constructor() {
    super();
    this.state = {
      isExpanded: false
    };
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
      this.props.onAdd({ id: 1, name: "New site", unsaved: true })
  }

  @bind
  expand() {
    this.setState({
      isExpanded: true
    });
  }

  @bind
  collapse() {
    this.setState({
      isExpanded: false
    });
  }

  @bind
  toggleExpanded() {
    this.setState({
      isExpanded: !this.state.isExpanded
    });
  }

  render(props, state) {
    // show nothing if there is only 1 site and no option to add additional sites
    if(!props.showAdd && props.sites.length == 1) {
        return '';
    }  

    // otherwise, render list of sites + add button
    return (
        <li style={{display: props.visible ? '' : 'none'}} class={`sites ${state.isExpanded ? 'expanded' : ''}`} onClick={this.toggleExpanded} onMouseEnter={this.expand} onMouseLeave={this.collapse}>
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

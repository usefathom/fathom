'use strict';

import { h, Component } from 'preact';
import { bind } from 'decko';

class SiteSwitcher extends Component {

  componentDidMount() {
     
  }

  componentWillUnmount() {
      
  }

  @bind 
  addSite() {
      this.props.onAdd({ id: 0, name: "New site"})
  }

  render(props, state) {
    return (
        <li class="sites">
            <a href="javascript:void(0)">Current site</a>
            <ul>
                <li class="add-new"><a href="javascript:void(0);" onClick={this.addSite}>+ Add another site</a></li>
            </ul>
        </li>
    )
  }
}

export default SiteSwitcher

'use strict';

import { h, Component } from 'preact';
import Client from '../lib/client.js';
import { bind } from 'decko';

class SiteSwitcher extends Component {

  constructor(props) {
    super(props)
    this.state = {}
  }

  componentDidMount() {
      this.fetchData();
  }

  componentWillUnmount() {
      
  }

  @bind
  fetchData() {
   
  }

  render(props, state) {
    return (
        <li class="sites">
            <a href="#">Current site</a>
            <ul>
                <li class="add-new"><a href="#">+ Add another site</a></li>
            </ul>
        </li>
    )
  }
}

export default SiteSwitcher

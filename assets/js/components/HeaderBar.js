'use strict'

import { h, render, Component } from 'preact';
import LogoutButton from '../components/LogoutButton.js';

class HeaderBar extends Component {
  render() {
    const rightContent = this.props.showLogout ? <LogoutButton onSuccess={this.props.onLogout} /> : '';

    return (
      <header class="header-bar cf">
        <div class="container">
          <h1 class="pull-left title">Ana <small class="subtitle">open web analytics</small></h1>
          <div class="pull-right">
              {rightContent}
          </div>
        </div>
      </header>
  )}
}

export default HeaderBar

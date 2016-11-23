'use strict';

import { h, render, Component } from 'preact';

class LogoutButton extends Component {

  constructor(props) {
    super(props)
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleSubmit(e) {
    e.preventDefault();

    fetch('/api/session', {
      method: "DELETE",
      credentials: 'include',
    }).then((r) => {
      if( r.status == 200 ) {
        this.props.onSuccess();
        console.log("No longer authenticated!");
      }
    });
  }

  render() {
    return (
      <a href="#" onClick={this.handleSubmit}>Sign out</a>
    )
  }
}

export default LogoutButton

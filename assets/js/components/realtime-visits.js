'use strict';

import React, { Component } from 'react'

class RealtimeVisitsCount extends React.Component {
  constructor(props) {
    super(props);
    this.state = { count: 0 };
    this.refresh() && window.setInterval(this.refresh.bind(this), 5000);
  }

  refresh() {
    return fetch('/api/visits/count/realtime')
      .then((r) => r.json())
      .then((data) => {
        this.setState({count: data});
    });
  }

  render() {
    let visitors = this.state.count > 1 ? 'visitors' : 'visitor';

    return (
      <div className="block">
        <span className="count">{this.state.count}</span> <span>{visitors} on the site right now.</span>
      </div>
    );
  }
}

export default RealtimeVisitsCount;

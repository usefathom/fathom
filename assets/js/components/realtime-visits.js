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
    return (
      <div>
        <h2>Real-time Visitors</h2>
        <span className="count">{this.state.count}</span>
      </div>
    );
  }
}

export default RealtimeVisitsCount;

'use strict'

import { h, render, Component } from 'preact';

class Notification extends Component {
  constructor(props) {
    super(props)

    this.state = {
      message: props.message,
      kind: props.kind || 'error'
    }
  }

  clearMessage() {
    this.setState({ message: '' })
  }

  componentWillReceiveProps(newProps) {
    if(newProps.message != this.state.message) {
      this.setState({ message: newProps.message, kind: newProps.kind || 'error' })
      window.setTimeout(this.clearMessage.bind(this), 5000)
    }
  }

  render() {
    if(this.state.message === '') {
      return ''
    }

    return (
      <div class={`notification`}>
        <div class={`notification-${this.state.kind}`}>
          {this.state.message}
        </div>
      </div>
  )}
}

export default Notification

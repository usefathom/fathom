'use strict'

import { h, render, Component } from 'preact';
import { bind } from 'decko';

class Notification extends Component {
  constructor(props) {
    super(props)

    this.state = {
      message: props.message,
      kind: props.kind || 'error'
    }
    this.timeout = 0
  }

  componentWillReceiveProps(newProps) {
    if(newProps.message != this.state.message) {
      this.setState({ message: newProps.message, kind: newProps.kind || 'error' })
      window.clearTimeout(this.timeout)
      this.timeout = window.setTimeout(this.clearMessage, 5000)
    }
  }

  @bind
  clearMessage() {
    this.setState({ message: '' })
  }

  render(props, state) {
    if(state.message === '') {
      return ''
    }

    return (
      <div class={`notification`}>
        <div class={`notification-${state.kind}`}>
          {state.message}
        </div>
      </div>
  )}
}

export default Notification

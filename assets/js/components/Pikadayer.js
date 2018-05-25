'use strict';

import Pikaday from 'pikaday';
import { h, Component } from 'preact';

class Pikadayer extends Component {
  componentDidMount() {
    this.pikaday = new Pikaday({ 
      field: this.base,
      onSelect: this.props.onSelect,
      position: 'bottom right',
   })
  }

  componentWillReceiveProps(newProps) {
    // make sure pikaday updates if we set a date using one of our presets
    if(this.pikaday && newProps.value !== this.props.value) {
      this.pikaday.setDate(newProps.value, true)
    }
  }

  componentWillUnmount() {
    this.pikaday.destroy()
  }

  render(props) {
    return <input value={props.value} />
  }
}

export default Pikadayer

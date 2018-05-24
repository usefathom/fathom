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
    if(newProps.value !== this.state.value) {
      this.state.value = newProps.value;
      this.pikaday.setDate(newProps.value, false);
    }
  }

  render(props) {
    return <input {...props} />
  }
}

export default Pikadayer

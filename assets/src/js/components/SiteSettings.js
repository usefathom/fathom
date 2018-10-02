'use strict';

import { h, Component } from 'preact';
import Client from '../lib/client.js';
import { bind } from 'decko';

class SiteSettings extends Component {

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

    @bind 
    onSubmit(evt) {
        console.log(evt)
    }

    render(props, state) {
        return (
            <div class="modal-wrap" style={"display: " + ( props.visible ? '' : 'none')}>
                <div class="modal">
                    <p>Update your site name or get your tracking code</p>
                    <form onSubmit={this.onSubmit}>
                        <fieldset>
                            <label>Site Name</label>
                            <input type="text" name="site-name" id="sitename" placeholder="" value="Paul Jarvis" />
                        </fieldset>

                        <fieldset>
                            <label>Add this code to your website</label>
                            <textarea onclick={(evt) => evt.target.select() } readonly="readonly">{`<!-- Fathom - simple website analytics - https://github.com/usefathom/fathom -->
<script>
(function(f, a, t, h, o, m){
	a[h]=a[h]||function(){
		(a[h].q=a[h].q||[]).push(arguments)
	};
	o=f.createElement('script'),
	m=f.getElementsByTagName('script')[0];
	o.async=1; o.src=t; o.id='fathom-script';
	m.parentNode.insertBefore(o,m)
})(document, window, '//stats.usefathom.com/tracker.js', 'fathom');
fathom('trackPageview');
</script>
<!-- / Fathom -->`}
                            </textarea>
                            <small><a href="#">Copy code</a></small>
                        </fieldset>

                        <fieldset>
                            <div class="half">
                                <div class="submit"><button type="submit">Update site name</button></div>
                                <div class="delete"><a href="#">Delete site</a></div>
                            </div>
                        </fieldset>
                    </form>
                </div>
            </div>
        )
    }
}

export default SiteSettings

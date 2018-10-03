'use strict';

import { h, Component } from 'preact';
import Client from '../lib/client.js';
import { bind } from 'decko';

class SiteSettings extends Component {

    constructor(props) {
        super(props)

        this.state = {
            copied: false,
        }
    }

    componentDidMount() {}
    componentWillUnmount() {}

    @bind
    copyToClipboard(evt) {
        this.textarea.select()
        document.execCommand('copy')
        this.setState({ copied: true })
        window.setTimeout(() => { this.setState({copied: false})}, 2400)
    }

    @bind 
    deleteSite(evt) {
        if(!confirm("Are you sure you want to delete this site? This action is irreversible - you will lose all the site's data too.")) {
            return;
        }

        Client.request(`/sites/${this.props.site.id}`, {
            method: "DELETE",
          }).then((r) => {
              console.log(r)
          })
    }

    @bind 
    onSubmit(evt) {
        evt.preventDefault();
       
        Client.request('sites', {
            method: "POST",
            data: {
              id: this.props.site.id,
              name: this.props.site.name,
            }
          }).then((r) => {
              console.log(r)
          })

    }

    @bind 
    onTextareaClick(evt) {
        evt.target.select()
    }

    @bind 
    maybeCloseModal(evt) {
        // don't close if click was inside the modal
        if ( evt.target.matches('.modal *, .modal')) {
            return;
        }

        this.props.onClose()
    }

    @bind 
    updateSiteName(evt) {
        this.props.site.name = evt.target.value;
    }

    render(props, state) {
        // TODO: Render different form for new sites vs. existing sites
        return (
        <div class="modal-wrap" style={"display: " + ( props.visible ? '' : 'none')} onClick={this.maybeCloseModal}>
            <div class="modal">
                <p>Update your site name or get your tracking code</p>
                <form onSubmit={this.onSubmit}>
                    <fieldset>
                        <label for="site-name">Site Name</label>
                        <input type="text" name="site-name" id="site-name" placeholder="" onChange={this.updateSiteName} value={props.site.name} />
                    </fieldset>

                    <fieldset>
                        <label>Add this code to your website</label>
                        <textarea ref={(el) => { this.textarea = el }} onClick={this.onTextareaClick} readonly="readonly">{`<!-- Fathom - simple website analytics - https://github.com/usefathom/fathom -->
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
                    <small><a href="javascript:void(0);" onClick={this.copyToClipboard}>{state.copied ? "Copied!" : "Copy code"}</a></small>
                </fieldset>

                <fieldset>
                    <div class="half">
                        <div class="submit"><button type="submit">Update site name</button></div>
                        <div class="delete"><a href="javascript:void(0);" onClick={this.deleteSite}>Delete site</a></div>
                    </div>
                </fieldset>
            </form>
        </div>
    </div>)
    }
}

export default SiteSettings

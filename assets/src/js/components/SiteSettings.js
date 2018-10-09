'use strict';

import { h, Component } from 'preact';
import Client from '../lib/client.js';
import { bind } from 'decko';

class SiteSettings extends Component {
    constructor(props) {
        super(props)

        this.state = {
            copied: false,
            updated: false,
        }
    }

    componentDidMount() {
        document.addEventListener('keydown', this.handleKeydownEvent);
    }
    componentWillUnmount() {
        document.removeEventListener('keydown', this.handleKeydownEvent)
    }

    @bind
    revertTemporaryState() {
        this.setState({
            copied: false, 
            updated: false
        })
    }

    @bind
    copyToClipboard(evt) {
        this.textarea.select()
        document.execCommand('copy')
        this.setState({ copied: true })
        window.setTimeout(this.revertTemporaryState, 2400)
    }

    @bind 
    deleteSite(evt) {
        if(!confirm("Are you sure you want to delete this site? This action is irreversible - you will lose all the site's data.")) {
            return;
        }

        let site = this.props.site;
        Client.request(`/sites/${site.id}`, {
            method: "DELETE",
          }).then((d) => {
              this.props.onDelete(site)
          })
    }

    @bind 
    onSubmit(evt) {
        evt.preventDefault();
        let site = this.props.site;
        let url = site.unsaved ? `/sites` : `/sites/${site.id}`

        Client.request(url, {
            method: "POST",
            data: {
                name: site.name,
            },
          }).then((site) => {
            this.setState({ updated: true})
            window.setTimeout(this.revertTemporaryState, 2400)

            site.unsaved = false
            this.props.onUpdate(site)
          })
    }

    @bind 
    handleTextareaClickEvent(evt) {
        evt.target.select()
    }

    @bind 
    handleClickEvent(evt) {
        // don't close if click was inside the modal
        if ( evt.target.matches('.modal *, .modal')) {
            return;
        }

        this.props.onClose()
    }

    @bind 
    handleKeydownEvent(evt) {
        // close modal when pressing ESC 
        if(evt.keyCode == 27) {
            this.props.onClose()
        }
    }

    @bind
    setTextarea(el) {
       this.textarea = el
    }

    @bind 
    updateSiteName(evt) {
        this.props.site.name = evt.target.value;
    }

    render(props, state) {
        return (
        <div class="modal-wrap" style={"display: " + ( props.visible ? '' : 'none')} onClick={this.handleClickEvent}>
            <div class="modal">
                <p>{props.site.unsaved ? 'Add a new site to track with Fathom' : 'Update your site name or get your tracking code'}</p>
                <form onSubmit={this.onSubmit}>
                    <fieldset>
                        <label for="site-name">Site name</label>
                        <input type="text" name="site-name" id="site-name" placeholder="" onChange={this.updateSiteName} value={props.site.name} />
                    </fieldset>

                    <fieldset style={props.site.unsaved ? 'display: none;' : ''}>
                        <label>Add this code to your website    <small class="right">(site ID = {props.site.trackingId})</small></label>
                        <textarea ref={this.setTextarea} onFocus={this.handleTextareaClickEvent} readonly="readonly">{`<!-- Fathom - simple website analytics - https://github.com/usefathom/fathom -->
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
fathom('set', 'siteId', '${props.site.trackingId}');
fathom('trackPageview');
</script>
<!-- / Fathom -->`}
                    </textarea>
                    <small><a href="javascript:void(0);" onClick={this.copyToClipboard}>{state.copied ? "Copied!" : "Copy code"}</a></small>
                </fieldset>

                <fieldset>
                    <div class="half">
                        <div class="submit"><button type="submit">{props.site.unsaved ? 'Create site' : 'Update site name'}</button> &nbsp; {state.updated ? 'Saved!' : ''}</div>
                        {props.site.unsaved ? '' : (<div class="delete"><a href="javascript:void(0);" onClick={this.deleteSite}>Delete site</a></div>)}
                    </div>
                </fieldset>
            </form>
        </div>
    </div>)
    }
}

export default SiteSettings

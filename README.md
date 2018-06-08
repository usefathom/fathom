Fathom - simple website analytics
==============================

[![Go Report Card](https://goreportcard.com/badge/github.com/usefathom/fathom)](https://goreportcard.com/report/github.com/usefathom/fathom)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/usefathom/fathom/master/LICENSE)

Fathom Analytics is a simpler and more privacy-focused alternative to Google Analytics.

Collecting information on the internet is important, but it’s broken. We’ve become complacent in trading information for free access to web services, and then complaining when those web services do crappy things with that data.

The problem is this: _if we aren’t paying for the product, we are the product_.

Google Analytics may give you free access to their services but in turn, they’re assembling data profiles on your website visitors, which they can then use for better targeting of advertisements across their network.

We need to stop giving away our data and our users' privacy for free access to a tool.

Fathom respects the privacy of your users and does not collect any personally identifiable information. All while giving you the information you need about your site, so you can make smarter decisions about your design and content.

![Screenshot of the Fathom dashboard](https://github.com/usefathom/fathom/raw/master/assets/src/img/fathom.jpg?v=7)

## Installation

For getting a development version of Fathom up & running, go through the following steps.

1. Ensure you have [Golang](https://golang.org/doc/install#install) installed properly
1. Get code: `git clone git@github.com:usefathom/fathom.git $GOPATH/src/github.com/usefathom/fathom` 
1. Compile into binary & prepare assets: `make build` 
1. (Optional) Set your [custom configuration values](https://github.com/usefathom/fathom/wiki/Configuration-file).
1. Register your user account: `fathom register --email=<email> --password=<password>`
1. Start the webserver: `fathom server` and then visit **http://localhost:8080** to access your analytics dashboard.

To install and run Fathom in production, [have a look at the installation instructions](https://github.com/usefathom/fathom/wiki/Installing-&-running-Fathom).

## Building with Docker

Ensure you have Docker installed and run `docker build -t fathom .`.
Run the container with `docker run -d -p 8080:8080 fathom`.

## Running with Docker

To run [our pre-built Docker image](https://hub.docker.com/r/usefathom/fathom/), run `docker run -d -p 8080:8080 usefathom/fathom:latest`

#### Tracking snippet 

To start tracking, include the following JavaScript on your site and replace `yourfathom.com` with the URL to your Fathom instance.

```html
<!-- Fathom - simple website analytics - https://github.com/usefathom/fathom -->
<script>
(function(f, a, t, h, o, m){
	a[h]=a[h]||function(){
		(a[h].q=a[h].q||[]).push(arguments)
	};
	o=f.createElement('script'),
	m=f.getElementsByTagName('script')[0];
	o.async=1; o.src=t; o.id='fathom-script';
	m.parentNode.insertBefore(o,m)
})(document, window, '//yourfathom.com/tracker.js', 'fathom');
fathom('trackPageview');
</script>
<!-- / Fathom -->
```

## Copyright and license

MIT licensed. Fathom and Fathom logo are trademarks of Fathom Analytics.

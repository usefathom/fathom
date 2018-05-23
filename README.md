Fathom - simple website analytics
==============================

[![Go Report Card](https://goreportcard.com/badge/github.com/usefathom/fathom)](https://goreportcard.com/report/github.com/usefathom/fathom)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/usefathom/fathom/master/LICENSE)


This is nowhere near being usable, let alone stable. Please treat as a proof of concept while we work on getting this to a stable state. **Do not run Fathom in production yet unless you like spending time on it.** Things will keep changing for the next few months.

![Screenshot of the Fathom dashboard](https://github.com/usefathom/fathom/blob/master/assets/img/fathom.jpg?v=6)

## Installation

For getting a development version of Fathom up & running, please go through the following steps.

1. Ensure you have [Golang](https://golang.org/doc/install#install) installed properly
1. Get code: `git clone git@github.com:usefathom/fathom.git $GOPATH/src/github.com/usefathom/fathom` 
1. Compile into binary & prepare assets: `make all` 
1. Register your user account: `fathom register --email=<email> --password=<password>`
1. Start the webserver: `fathom server` and then visit **http://localhost:8080** to access your analytics dashboard.

To start tracking, include the following JavaScript on your site and replace `yourfathom.com` with the URL to your Fathom instance (twice).

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

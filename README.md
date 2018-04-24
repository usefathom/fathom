Fathom - simple website analytics
==============================

[![Go Report Card](https://goreportcard.com/badge/github.com/usefathom/fathom)](https://goreportcard.com/report/github.com/usefathom/fathom)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/usefathom/fathom/master/LICENSE)


This is nowhere near being usable, let alone stable. Please treat as a proof of concept while we work on getting this to a stable state. **Do not run Fathom in production yet unless you like spending time on it.** Things will keep changing for the next few months.

![Screenshot of the Fathom dashboard](https://github.com/usefathom/fathom/raw/master/assets/dist/img/screenshot.png?v=6)

## Installation

For getting a development version of Fathom up & running, please go through the following steps.

1. get code: `go get -u github.com/usefathom/fathom` (or `git clone` repo into your `$GOPATH` )
1. run `npm install` (in code directory) to install all required dependencies
1. Rename `.env.example` to `.env` and set your database credentials.
1. Compile into binary: `make`
1. Create your user account: `fathom register <email> <password>`
1. Run default Gulp task to build static assets: `gulp`
1. Start the webserver: `fathom server --port=8080` & visit **localhost:8080** to access your analytics dashboard.

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
	o.async=1; o.src=t;
	m.parentNode.insertBefore(o,m)
})(document, window, '//yourfathom.com/tracker.js', 'fathom');
fathom('setTrackerUrl', '//yourfathom.com/collect');
fathom('trackPageview');
</script>
<!-- / Fathom -->
```

## License

MIT licensed.

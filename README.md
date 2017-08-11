Ana. Open Source Web Analytics.
==============================

[![Go Report Card](https://goreportcard.com/badge/github.com/dannyvankooten/ana)](https://goreportcard.com/report/github.com/dannyvankooten/ana)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/dannyvankooten/ana/master/LICENSE)


This is nowhere near being usable, let alone stable. Please treat as a proof of concept while we work on getting this to a stable state. **Do not run Ana in production yet unless you like spending time on it.** Things will keep changing for the next few months.

![Screenshot of the Ana dashboard](https://github.com/dannyvankooten/ana/raw/master/assets/img/screenshot.png?v=6)

## Installation

For getting a development version of Ana up & running, please go through the following steps.

1. get code: `go get -u github.com/dannyvankooten/ana` (or `git clone` repo into your `$GOPATH` )
1. run `npm install` (in code directory) to install all required dependencies
1. Rename `.env.example` to `.env` and set your database credentials.
1. Compile into binary: `make`
1. Create your user account: `ana register <email> <password>`
1. Run default Gulp task to build static assets: `gulp`
1. Start the webserver: `ana server --port=8080` & visit **localhost:8080** to access your analytics dashboard.

To start tracking, include the following JavaScript on your site and replace `ana.dev` with the URL to your Ana instance.

```html
<!-- Ana tracker -->
<script>
(function(d, w, u, o){
	w[o]=w[o]||function(){
		(w[o].q=w[o].q||[]).push(arguments)
	};
	a=d.createElement('script'),
	m=d.getElementsByTagName('script')[0];
	a.async=1; a.src=u;
	m.parentNode.insertBefore(a,m)
})(document, window, '//ana.dev/tracker.js', 'ana');
ana('setTrackerUrl', '//ana.dev/collect');
ana('trackPageview');
</script>
<!-- / Ana tracker -->
```

## License

MIT licensed.

Fathom Lite - simple website analytics
==============================
[![Go Report Card](https://goreportcard.com/badge/github.com/usefathom/fathom)](https://goreportcard.com/report/github.com/usefathom/fathom)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/usefathom/fathom/master/LICENSE)

Fathom Lite is a previous and open-source version of [Fathom Analytics](https://usefathom.com) (a paid, hosted [Google Analytics alternative](https://usefathom.com/google-analytics-alternative)). It was the very first version of our software, and has been downloaded millions of times!

While we are no longer adding features to this Lite version, we will be continuing to maintain it long-term and fix any bugs that come up.

![Screenshot of the Fathom dashboard](https://github.com/usefathom/fathom/raw/master/assets/src/img/fathom.jpg?v=7)

## Fathom Lite vs Fathom (hosted)

Today’s [Fathom Analytics](https://usefathom.com) is a hosted product with [simple pricing](https://usefathom.com/pricing) based on monthly pageviews. The same core capabilities are included on every plan: for example [API access](https://usefathom.com/api), up to 50 sites, custom events and ecommerce tracking, unlimited email reports and CSV exports, and [forever data retention](https://usefathom.com/features)—with [privacy-law compliance](https://usefathom.com/compliance) and no cookie banner required for analytics.

If you’d rather not run servers or maintenance yourself, try a [30-day free trial](https://usefathom.com/ref/GITHUB) (that link applies a **$10 credit** on your first invoice). Browse [all features](https://usefathom.com/features) or the [live demo](https://app.usefathom.com/demo).

![Screenshot of the Fathom Analytics Dashboard](https://usefathom.com/assets/images/fathom-screenshot.png)

| Feature | Fathom Lite | Fathom (hosted) |
|---------|-------------|-----------------|
| Fully managed | ✗ Self-hosted; you run servers and updates | ✓ Managed for you; [pay per pageviews](https://usefathom.com/pricing) |
| Cookie-free (no analytics banner) | ✗ Uses cookies | ✓ [Cookie-free](https://usefathom.com/features) tracking |
| Current dashboard (real-time, live visitors, filters, details) | ✗ Older Lite UI | ✓ Full dashboard — see [features](https://usefathom.com/features) |
| API | ✗ | ✓ [API](https://usefathom.com/api) on all plans |
| Custom events, ecommerce, UTMs | ✗ | ✓ Events, revenue, campaigns |
| [EU isolation](https://usefathom.com/features/eu-isolation) & [custom domains](https://usefathom.com/features/custom-domains) | ✗ | ✓ EU routing, first-party script domains |
| [GA import](https://usefathom.com/features/ga-importer), email reports, CSV export, [shared dashboards](https://usefathom.com/features) | ✗ | ✓ Unlimited reports & exports |
| Many sites per account | ✗ Typical single install | ✓ Up to 50 sites (more available) |
| Email support & SLA | ✗ Community / as-is | ✓ Support on every plan |
| Global CDN, scaling, backups | ✗ Your responsibility | ✓ Included |
| Active feature development | ✗ Bugfixes / maintenance | ✓ Ongoing — [trial](https://usefathom.com/ref/GITHUB) · [sign up](https://app.usefathom.com/register) |


## Installation


### Production

You can install Fathom on your server by following [our simple instructions](docs/Installation%20instructions.md).

### Development

For getting a development version of Fathom up & running, go through the following steps.

1. Ensure you have [Go](https://golang.org/doc/install#install) and [NPM](https://www.npmjs.com) installed
1. Download the code: `git clone https://github.com/usefathom/fathom.git $GOPATH/src/github.com/usefathom/fathom`
1. Compile the project into an executable: `make build`
1. (Optional) Set [custom configuration values](docs/Configuration.md)
1. (Required) Register a user account: `./fathom user add --email=<email> --password=<password>`
1. Start the webserver: `./fathom server` and then visit **http://localhost:8080** to access your analytics dashboard

## Docker

### Building

Ensure you have Docker installed and run `docker build -t fathom .`.
Run the container with `docker run -d -p 8080:8080 fathom`.

### Running

To run [our pre-built Docker image](https://hub.docker.com/r/usefathom/fathom/), run `docker run -d -p 8080:8080 usefathom/fathom:latest`

## Tracking snippet

To start tracking, create a site in your Fathom dashboard and copy the tracking snippet to the website(s) you want to track.

### Content Security Policy

If you use a [Content Security Policy (CSP)](https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP) to specify security policies for your website, Fathom requires the following CSP directives (replace `yourfathom.com` with the URL to your Fathom instance):

```
script-src: yourfathom.com;
img-src: yourfathom.com;
```

## Copyright and license

MIT licensed. Fathom and Fathom logo are trademarks of Fathom Analytics.

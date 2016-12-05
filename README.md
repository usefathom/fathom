Ana. Open Source Web Analytics.
==============================

This is nowhere near being usable, let alone stable. Treat as a proof of concept.

![Screenshot of the Ana dashboard](https://github.com/dannyvankooten/ana/raw/master/assets/img/screenshot.png?v=6)

## Usage

```html
<!-- Ana tracker -->
<script>
(function(d, w, u, o){
	w[o]=w[o]||function(){
		(w[o].q=w[o].q||[]).push(arguments)
	};
	a=d.createElement('script'),
	m=d.getElementsByTagName('script')[0];
	a.async=1;
	a.src=u;
	m.parentNode.insertBefore(a,m)
})(document, window, '//ana.dev/tracker.js', 'ana');
ana('setTrackerUrl', '//ana.dev/collect');
ana('trackPageview');
</script>
<!-- / Ana tracker -->
```

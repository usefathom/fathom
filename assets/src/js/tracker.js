'use strict';

var queue = window.ana.q || [];
var trackerUrl = '//ana.dev/collect';
var commands = {
  "trackPageview": trackPageview,
  "setTrackerUrl": setTrackerUrl,
};

// convert object to query string
function stringifyObject(json) {
  var keys = Object.keys(json);

  // omit empty
  keys = keys.filter(function(k) {
    return json[k].length > 0;
  });

  return '?' +
      keys.map(function(k) {
          return encodeURIComponent(k) + '=' +
              encodeURIComponent(json[k]);
      }).join('&');
}

function setTrackerUrl(v) {
  trackerUrl = v;
}

function trackPageview() {
  // Respect "Do Not Track" requests
  if( navigator.DoNotTrack === "1" ) {
    return;
  }

  // get the path or canonical
  var path = location.pathname + location.search;
  var canonical = document.querySelector('link[rel="canonical"]');
  if(canonical && canonical.href) {
    path = canonical.href.substring(canonical.href.indexOf('/', 7)) || '/';
  }

  var d = {
    h: location.hostname,
    t: document.title,
    l: navigator.language,
    p: path,
    sr: screen.width + "x" + screen.height,
    t: document.title,
    ru: document.referrer,
    rk: ""
  };

  var i = document.createElement('img');
  i.src = trackerUrl + stringifyObject(d);
  document.body.appendChild(i);
}

// override global ana object
window.ana = function() {
  var args = [].slice.call(arguments);
  var c = args.shift();
  commands[c].apply(this, args);
};

// process existing queue
queue.forEach(function(i) {
  ana.apply(this, i);
});

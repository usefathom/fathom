'use strict';

var queue = window.fathom.q || [];
var trackerUrl = '';
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
  if(trackerUrl === '') {
    console.error('Fathom: invalid tracker URL');
    return;
  }

  // Respect "Do Not Track" requests
  if(navigator.DoNotTrack === "1") {
    return;
  }

  // get the path or canonical
  var path = location.pathname + location.search;
  var canonical = document.querySelector('link[rel="canonical"][href]');
  if(canonical) {
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

// override global fathom object
window.fathom = function() {
  var args = [].slice.call(arguments);
  var c = args.shift();
  commands[c].apply(this, args);
};

// process existing queue
queue.forEach(function(i) {
  fathom.apply(this, i);
});

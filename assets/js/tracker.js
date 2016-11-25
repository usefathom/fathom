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
  if( navigator.DonotTrack == 1 ) {
    return;
  }

  var d = {
    l: navigator.language,
    p: location.pathname + location.search,
    sr: screen.width + "x" + screen.height,
    t: document.title,
    r: document.referrer
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
queue.map((i) => ana.apply(this, i));

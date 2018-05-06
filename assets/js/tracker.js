'use strict';

var cookies = require('cookies-js');
var queue = window.fathom.q || [];
var trackerUrl = '';
var commands = {
  "trackPageview": trackPageview,
  "setTrackerUrl": setTrackerUrl,
};

// convert object to query string
function stringifyObject(json) {
  var keys = Object.keys(json);

  return '?' +
      keys.map(function(k) {
          return encodeURIComponent(k) + '=' +
              encodeURIComponent(json[k]);
      }).join('&');
}

function generateKey() {
  var s = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
  return Array(16).join().split(',').map(function() { return s.charAt(Math.floor(Math.random() * s.length)); }).join('');
}

function getData() {
  var data = Cookies.get('_fathom');

  if(data) {
    try{
      data = JSON.parse(data);
      return data;
    } catch(e) {
      console.error(e);
    }
  }

  return {
    sid: generateKey(),
    isNew: true, 
    pagesViewed: [],
  }
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
  if('doNotTrack' in navigator && navigator.doNotTrack === "1") {
    return;
  }

  // get the path or canonical
  var path = location.pathname + location.search;

  // parse path from canonical, if page has one
  var canonical = document.querySelector('link[rel="canonical"][href]');
  if(canonical) {
    var a = document.createElement('a');
    a.href = canonical.href;
    path = a.pathname;
  }

  let referrer = '';
  if(document.referrer.indexOf(location.hostname) < 0) {
    referrer = document.referrer;
  }

  let data = getData();

  var d = {
    sid: data.sid,
    p: path,
    t: document.title,
    r: referrer,
    scheme: location.protocol.substring(0, location.protocol.length - 1),
    u: data.pagesViewed.indexOf(path) == -1 ? 1 : 0,
    b: data.isNew ? 1 : 0, // because only new visitors can bounce. we update this server-side.
    n: data.isNew ? 1 : 0, 
  };

  var i = document.createElement('img');
  i.src = trackerUrl + stringifyObject(d);
  i.addEventListener('load', function() {
    data.pagesViewed.push(path);
    data.isNew = false;
    Cookies.set('_fathom', JSON.stringify(data), { expires: 60 * 60 * 24});
  });
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

'use strict';

/* cookies.js - https://github.com/ScottHamper/Cookies */
(function(d,f){"use strict";var h=function(d){if("object"!==typeof d.document)throw Error("Cookies.js requires a `window` with a `document` object");var b=function(a,e,c){return 1===arguments.length?b.get(a):b.set(a,e,c)};b._document=d.document;b._cacheKeyPrefix="cookey.";b._maxExpireDate=new Date("Fri, 31 Dec 9999 23:59:59 UTC");b.defaults={path:"/",secure:!1};b.get=function(a){b._cachedDocumentCookie!==b._document.cookie&&b._renewCache();a=b._cache[b._cacheKeyPrefix+a];return a===f?f:decodeURIComponent(a)};
b.set=function(a,e,c){c=b._getExtendedOptions(c);c.expires=b._getExpiresDate(e===f?-1:c.expires);b._document.cookie=b._generateCookieString(a,e,c);return b};b.expire=function(a,e){return b.set(a,f,e)};b._getExtendedOptions=function(a){return{path:a&&a.path||b.defaults.path,domain:a&&a.domain||b.defaults.domain,expires:a&&a.expires||b.defaults.expires,secure:a&&a.secure!==f?a.secure:b.defaults.secure}};b._isValidDate=function(a){return"[object Date]"===Object.prototype.toString.call(a)&&!isNaN(a.getTime())};
b._getExpiresDate=function(a,e){e=e||new Date;"number"===typeof a?a=Infinity===a?b._maxExpireDate:new Date(e.getTime()+1E3*a):"string"===typeof a&&(a=new Date(a));if(a&&!b._isValidDate(a))throw Error("`expires` parameter cannot be converted to a valid Date instance");return a};b._generateCookieString=function(a,b,c){a=a.replace(/[^#$&+\^`|]/g,encodeURIComponent);a=a.replace(/\(/g,"%28").replace(/\)/g,"%29");b=(b+"").replace(/[^!#$&-+\--:<-\[\]-~]/g,encodeURIComponent);c=c||{};a=a+"="+b+(c.path?";path="+
c.path:"");a+=c.domain?";domain="+c.domain:"";a+=c.expires?";expires="+c.expires.toUTCString():"";return a+=c.secure?";secure":""};b._getCacheFromString=function(a){var e={};a=a?a.split("; "):[];for(var c=0;c<a.length;c++){var d=b._getKeyValuePairFromCookieString(a[c]);e[b._cacheKeyPrefix+d.key]===f&&(e[b._cacheKeyPrefix+d.key]=d.value)}return e};b._getKeyValuePairFromCookieString=function(a){var b=a.indexOf("="),b=0>b?a.length:b,c=a.substr(0,b),d;try{d=decodeURIComponent(c)}catch(k){console&&"function"===
typeof console.error&&console.error('Could not decode cookie with key "'+c+'"',k)}return{key:d,value:a.substr(b+1)}};b._renewCache=function(){b._cache=b._getCacheFromString(b._document.cookie);b._cachedDocumentCookie=b._document.cookie};b._areEnabled=function(){var a="1"===b.set("cookies.js",1).get("cookies.js");b.expire("cookies.js");return a};b.enabled=b._areEnabled();return b},g=d&&"object"===typeof d.document?h(d):h;"function"===typeof define&&define.amd?define(function(){return g}):"object"===
typeof exports?("object"===typeof module&&"object"===typeof module.exports&&(exports=module.exports=g),exports.Cookies=g):d.Cookies=g})("undefined"===typeof window?this:window);

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

function generateKey() {
  var s = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
  return Array(16).join().split(',').map(function() { return s.charAt(Math.floor(Math.random() * s.length)); }).join('');
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

  // set cookie with random visitor key, valid for 24 hours.
  // visitor can delete this cookie in order to be forgotten
  if(!Cookies.get('_fathom')) {
    Cookies.set('_fathom', generateKey(), { expires: 60 * 60 * 24 });
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

  var d = {
    vk: Cookies.get('_fathom'),
    h: location.hostname,
    t: document.title,
    l: navigator.language,
    p: path,
    sr: screen.width + "x" + screen.height,
    t: document.title,
    ru: document.referrer,
    rk: "",
    scheme: location.protocol.substring(0, location.protocol.length - 1),
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

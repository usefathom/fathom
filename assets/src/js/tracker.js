(function() { 
  'use strict';

  let queue = window.fathom.q || [];
  let config = {
    'siteId': '',
    'trackerUrl': '',
  };
  const commands = {
    "set": set,
    "trackPageview": trackPageview,
    "setTrackerUrl": setTrackerUrl,
  };

  function set(key, value) {
    config[key] = value;
  }

  function setTrackerUrl(value) {
    return set("trackerUrl", value);
  }

  // convert object to query string
  function stringifyObject(obj) {
    var keys = Object.keys(obj);

    return '?' +
        keys.map(function(k) {
            return encodeURIComponent(k) + '=' + encodeURIComponent(obj[k]);
        }).join('&');
  }

  function getCookie(name) {
    var cookies = document.cookie ? document.cookie.split('; ') : [];
    
    for (var i = 0; i < cookies.length; i++) {
      var parts = cookies[i].split('=');
      if (decodeURIComponent(parts[0]) !== name) {
        continue;
      }

      var cookie = parts.slice(1).join('=');
      return decodeURIComponent(cookie);
    }

    return '';
  }

  function setCookie(name, data, args) {
    name = encodeURIComponent(name);
    data = encodeURIComponent(String(data));

    var str = name + '=' + data;

    if(args.path) {
      str += ';path=' + args.path;
    }
    if (args.expires) {
      str += ';expires='+args.expires.toUTCString();
    }

    document.cookie = str + ';SameSite=None;Secure';
  }

  function newVisitorData() {
    return {
      isNewVisitor: true, 
      isNewSession: true,
      pagesViewed: [],
      previousPageviewId: '',
      lastSeen: +new Date(),
    }
  }

  function getData() {
    let thirtyMinsAgo = new Date();
    thirtyMinsAgo.setMinutes(thirtyMinsAgo.getMinutes() - 30);

    let data = getCookie('_fathom');
    if(! data) {
      return newVisitorData();
    }

    try{
      data = JSON.parse(data);
    } catch(e) {
      console.error(e);
      return newVisitorData();
    }

    if(data.lastSeen < (+thirtyMinsAgo)) {
      data.isNewSession = true;
    }

    return data;  
  }

  function findTrackerUrl() {
    const el = document.getElementById('fathom-script')
    return el ? el.src.replace('tracker.js', 'collect') : '';
  }

  function trackPageview(vars) { 
    vars = vars || {};

    // Respect "Do Not Track" requests
    if('doNotTrack' in navigator && navigator.doNotTrack === "1") {
      return;
    }

    // ignore prerendered pages
    if( 'visibilityState' in document && document.visibilityState === 'prerender' ) {
      return;
    }

    // if <body> did not load yet, try again at dom ready event
    if( document.body === null ) {
      document.addEventListener("DOMContentLoaded", () => {
        trackPageview(vars);
      })
      return;
    }

    //  parse request, use canonical if there is one
    let req = window.location;

    // do not track if not served over HTTP or HTTPS (eg from local filesystem) and we're not in an Electron app
    if(req.host === '' && navigator.userAgent.indexOf("Electron") < 0) {
      return;
    }

    // find canonical URL
    let canonical = document.querySelector('link[rel="canonical"][href]');
    if(canonical) {
      let a = document.createElement('a');
      a.href = canonical.href;

      // use parsed canonical as location object
      req = a;
    }
    
    let path = vars.path || ( req.pathname + req.search );
    if(!path) {
      path = '/';
    }

    // determine hostname
    let hostname = vars.hostname || ( req.protocol + "//" + req.hostname );

    // only set referrer if not internal
    let referrer = vars.referrer || '';
    if(document.referrer.indexOf(hostname) < 0) {
      referrer = document.referrer;
    }

    let data = getData();
    const d = {
      pid: data.previousPageviewId || '',
      p: path,
      h: hostname,
      r: referrer,
      u: data.pagesViewed.indexOf(path) == -1 ? 1 : 0,
      nv: data.isNewVisitor ? 1 : 0, 
      ns: data.isNewSession ? 1 : 0,
      sid: config.siteId,
    };

    let url = config.trackerUrl || findTrackerUrl()
    let img = document.createElement('img');
    img.setAttribute('alt', '');
    img.setAttribute('aria-hidden', 'true');
    img.src = url + stringifyObject(d);
    img.addEventListener('load', function() {
      let midnight = new Date();
      midnight.setHours(24); midnight.setMinutes(0); midnight.setSeconds(0);

      // update data in cookie
      if( data.pagesViewed.indexOf(path) == -1 ) {
        data.pagesViewed.push(path);
      }
      data.previousPageviewId = d.id;
      data.isNewVisitor = false;
      data.isNewSession = false;
      data.lastSeen = +new Date();
      setCookie('_fathom', JSON.stringify(data), { expires: midnight, path: '/' });

      // remove tracking img from DOM
      document.body.removeChild(img)
    });
    
    // in case img.onload never fires, remove img after 1s & reset src attribute to cancel request
    window.setTimeout(() => { 
      if(!img.parentNode) {
        return;
      }

      img.src = ''; 
      document.body.removeChild(img)
    }, 1000);

    // add to DOM to fire request
    document.body.appendChild(img);  
  }

  // override global fathom object
  window.fathom = function() {
    var args = [].slice.call(arguments);
    var c = args.shift();
    commands[c].apply(this, args);
  };

  // process existing queue
  queue.forEach((i) => fathom.apply(this, i));
})()

'use strict';

import * as cookies from 'cookies-js';
import * as util from './lib/util.js';

let queue = window.fathom.q || [];
let trackerUrl = '';

const commands = {
  "trackPageview": trackPageview,
  "setTrackerUrl": setTrackerUrl,
};

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

  let data = cookies.get('_fathom');
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

function setTrackerUrl(v) {
  trackerUrl = v;
}

function trackPageview() {
  if(trackerUrl === '') {
    trackerUrl = findTrackerUrl();
  }

  // Respect "Do Not Track" requests
  if('doNotTrack' in navigator && navigator.doNotTrack === "1") {
    return;
  }

  // ignore prerendered pages
  if( 'visibilityState' in document && document.visibilityState === 'prerender' ) {
    return;
  }

  let req = document.location;

  // parse canonical, if page has one
  let canonical = document.querySelector('link[rel="canonical"][href]');
  if(canonical) {
    let a = document.createElement('a');
    a.href = canonical.href;

    // use parsed canonical as location object
    req = a;
  }

  if (URL) {
    req = new URL(req);
    let path = req.pathname + req.search;
    let hostname = req.origin;
  }
  // if parsing path failed, default to main page and origin
  if(!path) {
    let path = '/';
  }
  if(!hostname) {
    let hostname = document.location.origin;
  }

  // only set referrer if not internal
  let referrer = '';
  if(document.referrer.indexOf(document.location.origin) < 0) {
    referrer = document.referrer;
  }

  let data = getData();
  const d = {
    id: util.randomString(20),
    pid: data.previousPageviewId || '',
    p: path,
    h: hostname,
    r: referrer,
    u: data.pagesViewed.indexOf(path) == -1 ? 1 : 0,
    nv: data.isNewVisitor ? 1 : 0, 
    ns: data.isNewSession ? 1 : 0,
  };

  let i = document.createElement('img');
  i.src = trackerUrl + util.stringifyObject(d);
  i.addEventListener('load', function() {
    let now = new Date();
    let midnight = new Date(Date.UTC(now.getFullYear(), now.getMonth(), now.getDate(), 24, 0, 0));
    let expires = Math.round((midnight - now) / 1000);

    // update data in cookie
    if( data.pagesViewed.indexOf(path) == -1 ) {
      data.pagesViewed.push(path);
    }
    data.previousPageviewId = d.id;
    data.isNewVisitor = false;
    data.isNewSession = false;
    data.lastSeen = +new Date();
    cookies.set('_fathom', JSON.stringify(data), { 'expires': expires });
  });
  document.body.appendChild(i);
  window.setTimeout(() => { document.body.removeChild(i)}, 1000);
}

// override global fathom object
window.fathom = function() {
  var args = [].slice.call(arguments);
  var c = args.shift();
  commands[c].apply(this, args);
};

// process existing queue
queue.forEach((i) => fathom.apply(this, i));

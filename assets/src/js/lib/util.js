'use strict';

// convert object to query string
function stringifyObject(json) {
  var keys = Object.keys(json);

  return '?' +
      keys.map(function(k) {
          return encodeURIComponent(k) + '=' +
              encodeURIComponent(json[k]);
      }).join('&');
}

function randomString(n) {
  var s = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
  return Array(n).join().split(',').map(() => s.charAt(Math.floor(Math.random() * s.length))).join('');
}

export { 
   randomString, 
   stringifyObject
}

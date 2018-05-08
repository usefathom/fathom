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

function generateKey() {
  var s = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
  return Array(16).join().split(',').map(function() { return s.charAt(Math.floor(Math.random() * s.length)); }).join('');
}

export { 
   generateKey, 
   stringifyObject
}

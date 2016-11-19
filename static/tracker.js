function jsonToQueryString(json) {
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

// abort hit if Do Not Track is enabled.
// if( navigator.DonotTrack == 1 ) {
//   return;
// }

var i = document.createElement('img');
var d = {
  l: navigator.language,
  p: location.pathname + location.search,
  sr: screen.width + "x" + screen.height,
  t: document.title,
  r: document.referrer
};


i.src = 'http://localhost:8080/collect' + jsonToQueryString(d);

document.body.appendChild(i);

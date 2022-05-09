'use strict';

// convert object to query string
function stringifyObject(json) {
    var keys = Object.keys(json);

    return '?' +
        keys.map(function (k) {
            return encodeURIComponent(k) + '=' +
                encodeURIComponent(json[k]);
        }).join('&');
}

function randomString(n) {
    var s = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
    return Array(n).join().split(',').map(() => s.charAt(Math.floor(Math.random() * s.length))).join('');
}

function hashParams() {
    var params = {},
        match,
        matches = window.location.hash.substring(2).split("&");

    for (var i = 0; i < matches.length; i++) {
        match = matches[i].split('=')
        params[match[0]] = decodeURIComponent(match[1]);
    }

    return params;
}

export {
    randomString,
    stringifyObject,
    hashParams
}

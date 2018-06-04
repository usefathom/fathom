'use strict';

function formatPretty(num) {
  let M = 1000000, K = 1000;

  if (num >= M) {
    return (num / M).toFixed(num > (M * 10) ? 1 : 2).replace(/\.0$/, '') + 'M';
  }

  if (num >= (K * 10)) {
    return (num / K).toFixed(num > (K*100) ? 0 : 1).replace(/\.0$/, '') + 'K';
  }

  return formatWithComma(num);
}

function formatWithComma(nStr) {
	nStr += '';

  if(nStr.length < 4 ) {
    return nStr;
  }

  var	x = nStr.split('.');
	var x1 = x[0];
	var x2 = x.length > 1 ? '.' + x[1] : '';
	var rgx = /(\d+)(\d{3})/;
	while (rgx.test(x1)) {
		x1 = x1.replace(rgx, '$1' + ',' + '$2');
	}
	return x1 + x2;
}

function formatDuration(seconds) {
  seconds = Math.round(seconds);
  var date = new Date(null);
  date.setSeconds(seconds); // specify value for SECONDS here
  return date.toISOString().substr(14, 5);
}

function formatPercentage(p) {
   return Math.round(p*100) + "%";
}

export { 
  formatPretty,
  formatWithComma, 
  formatDuration, 
  formatPercentage 
}

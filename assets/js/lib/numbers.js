'use strict';

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
   var date = new Date(null);
   date.setSeconds(seconds); // specify value for SECONDS here
   return date.toISOString().substr(14, 5);
}

export { formatWithComma, formatDuration }

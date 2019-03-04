 

var debug = false;
 
var authWindow;

var app;
 
var initApp = function( sid) {
	if (app!=null) { return; }
	app = new App( );
	app.init( document.getElementById('canvas'),sid )
  
}

var forceInit = function( sid ) {
	initApp( sid )
	document.getElementById('unsupported-browser').style.display = "none";
	return false;
}



 
function startWs( sid ) {

	// 开始启动
	if(Modernizr.canvas && Modernizr.websockets) {
		initApp( sid );
	} else {
		document.getElementById('unsupported-browser').style.display = "block";
		//document.getElementById('force-init-button').addEventListener('click', forceInit, false);
	}

	$('a[rel=external]').click(function(e) {
		e.preventDefault();
		window.open($(this).attr('href'));
	});
}

document.body.onselectstart = function() { return false; }

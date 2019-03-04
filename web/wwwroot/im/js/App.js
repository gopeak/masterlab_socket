
var App = function( aCanvas) {
	var app = this;

	var
			canvas,
			context,
			webSocket,
			webSocketService,
			messageQuota = 5,
			sid
	;

	app.update = function() {

	};

	app.draw = function() {

	};

	app.onSocketOpen = function(e) {
		var sendObj = {
				type: 'auth'
			};
		// 认证请求
		 app.authorize( GlobalToken, GlobalSid )

	};

	app.onSocketClose = function(e) {

		webSocketService.connectionClosed();
	};

	app.onSocketMessage = function(e) {

		console.log( e.data )
		data_json = JSON.parse( e.data )

		webSocketService.processMessage(data_json);

	};

	app.sendMessage = function( msg ) {

	    webSocketService.sendMessage( msg  );

	}
    app.pushMessage = function( from_sid,from_info,to_sid,msg ) {

        var sendObj = {
            sid: to_sid,
            from_info:from_info,
            msg: msg,
        };
        webSocketService.pushMessage( from_sid, sendObj  );

    }
	app.pushGroupMessage = function( from_sid,from_info,area_id,msg ) {

		var sendObj = {
			area_id: area_id,
			content: msg,
            from_info:from_info
		};
		webSocketService.pushGroupMessage( from_sid, sendObj  );

	}

	app.authorize = function(token,sid) {
		webSocketService.authorize(token,sid);
	}

	app.init = function(aCanvas ,sid) {
		canvas = aCanvas;
		context = canvas.getContext('2d');

		webSocket 				= new WebSocket( 'ws://'+document.domain+':9898/ws' );
		webSocket.onopen 		= app.onSocketOpen;
		webSocket.onclose		= app.onSocketClose;
		webSocket.onmessage 	= app.onSocketMessage;

		webSocketService		= new WebSocketService( webSocket );

	}

}

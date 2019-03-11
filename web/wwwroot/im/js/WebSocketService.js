var WebSocketService = function( webSocket) {
	var webSocketService = this;
	
	var webSocket = webSocket;

    var  TypeReq  	 = "1"
    var  TypeResp  	 = "2"
    var  TypeBroatcast  = "3"
    var  TypePush 	 = "4"
    var  TypeError  	 = "5"
    var  TypeReply	 = "6"
	
	this.hasConnection = false;
	
	this.welcomeHandler = function(json_obj) {

        webSocketService.hasConnection = true;
        console.log("welcomeHandler:",json_obj);

        webSocketService.subscripeGroup()
    };


    this.failedHandler = function(json_obj) {
        webSocketService.hasConnection = true;
        console.log("failedHandler:");
        console.log(json_obj);
        console.log("加入失败!")

    };

    this.errorHandler = function(json_obj) {
        console.log("errorHandler:");
        console.log("服务器返回错误:")
        console.log(json_obj);
    };

    this.subscripeGroupfailedHandler = function(json_obj) {
        console.log(json_obj);
        console.log("加入群组失败!")

    };

    this.subscripeGroupHandler = function(json_obj) {
        console.log(json_obj);
        console.log("加入群组成功!")

    };

    this.failedHandler = function(json_obj) {
        webSocketService.hasConnection = true;
        console.log("failedHandler:");
        console.log(json_obj);

    };
	
	this.updateHandler = function(json_obj) {
		var newtp = false;
		//console.log( "updateHandler:" );
		 console.log( json_obj );
 
	}
	
	this.pushHandler = function(json_obj ) {
		console.log( "messageHandler:" );

		data  = json_obj.data
        console.log( data);
        from_info = data.from_info
        var from_sid = data.sid
        /*
        for(var i=0; i<GlobalContacts.length; i++)
        {
            if  (GlobalContacts[i].sid==from_sid){
                from_info = GlobalContacts[i];
                break;
            }
        }*/

        obj = {
            username:from_info.username
            ,avatar: from_info.avatar
            ,id: from_info.id
            ,type: "friend"
			,mine:false
            ,content: data.msg
        }
        console.log( " layim.getMessage(obj):" );
        console.log( obj );
        layui.use('layim', function(layim){
            layim.getMessage(obj);
        });
		
	}

    this.broatcastHandler = function( json_obj ) {

		data  = JSON.parse(json_obj.data)
        console.log(  data );
        from_info = data.from_info
        console.log(  from_info );
        group_id = ""
        for(var i=0; i<GlobalGroups.length; i++)
        {
            if  (GlobalGroups[i].channel_id==data.area_id){
                group_id = GlobalGroups[i].id;
                break;
            }
        }

        var obj = {
            username:from_info.username
            ,avatar: from_info.avatar
            ,id: group_id
            ,fromid:from_info.id
			,mine:false
            ,type: "group"
            ,content: data.content
        }
        console.log( "messageGroupHandler obj:" );
        console.log( obj );

        layui.use('layim', function(layim){
            layim.getMessage(obj);
        });

    }
	
	this.closedHandler = function(json_obj) {

	}
	
	this.redirectHandler = function( json_obj ) {

		data = json_obj.data
		if (data.url) {
			if (authWindow) {
				authWindow.document.location = data.url;
			} else {
				document.location = data.url;
			}			
		}
	}
	
	this.noneHandler = function(json_obj) {
		 
	}
	
	this.processMessage = function( json_obj ) {
	    console.log("processMessage:");
        var fn

        if( json_obj.type==TypeError ) {
            fn = webSocketService[ 'errorHandler'];
            if (fn) {
                fn(json_obj);
            }
            return
        }

        if( json_obj.type==TypePush ) {
            fn = webSocketService[ 'pushHandler'];
            if (fn) {
                fn(json_obj);
            }
            return
        }
        if( json_obj.type==TypeBroatcast ) {
            fn = webSocketService[ 'broatcastHandler'];
            if (fn) {
                fn(json_obj);
            }
            return
        }

        if (typeof(json_obj.data) == "string") {
            try{
                json_obj.data = JSON.parse(json_obj.data)
            }catch(err){

            }
        }

        if (typeof(json_obj.data.type) == "undefined") {
             return
        }
        fn = webSocketService[json_obj.data.type + 'Handler'];
		if (fn) {
			fn(json_obj);
		}
	}
	
	this.connectionClosed = function() {
		webSocketService.hasConnection = false;
		 
	};
	
	this.wrapReqMessage = function( _cmd,sid,reqid,msg ){
	    // { "header":{ "cmd":"", "seq_id":0,  "sid":"" , "token":"", "version":"1.0" ,"gzip":true}  , "type":"req", "data":{}  }
		var req_obj = {
            header: {
				cmd:_cmd,
				seq_id:reqid,
				sid:sid, 
				token:GlobalToken,
				version:"1.0",
				no_resp:false,
				gzip:false
			},
            type:TypeReq,
            data: msg,
        };
        console.log( req_obj );
		return  JSON.stringify(req_obj) 

	}

 
	this.wrapPushMessage = function( sid,msg ){
		//  { "header":{ "cmd":"", "seq_id":0,  "sid":"" , "token":"", "version":"1.0" ,"gzip":true}  , "type":"req", "data":{}  }
		var req_obj = {
            header: {
				cmd:"PushMessage",
				seq_id:0,
				sid:sid, 
				token:GlobalToken,
				version:"1.0",
				no_resp:true,
				gzip:false
			},
            type:TypePush,
            data: msg,
        };
		return  JSON.stringify(req_obj)
	}
	 

	this.wrapPushGroupMessage = function( sid,msg ){
		//  { "header":{ "cmd":"", "seq_id":0,  "sid":"" , "token":"", "version":"1.0" ,"gzip":true}  , "type":"req", "data":{}  }
		var req_obj = {
            header: {
				cmd:"PushGroupMessage",
				seq_id:0,
				sid:sid, 
				token:GlobalToken,
				version:"1.0",
				no_resp:true,
				gzip:false
			},
            type:TypeBroatcast,
            data: msg,
        };
        return  JSON.stringify(req_obj)
	}

	this.sendMessage = function( sid, msg  ) {
		var sendObj = {
			type: 'message',
			message: msg,
			id:sid
		};
        str = this.wrapReqMessage( 'Message',sid,0,sendObj)
		webSocket.send(str);
	}
	 

	this.pushMessage = function( sid, msg  ) {
		console.log("pushMessage:");
        console.log( sid );
        console.log( msg );
		str = this.wrapPushMessage( sid,msg)
		webSocket.send(str);
	}
	this.pushGroupMessage = function( sid, msg  ) {

		str = this.wrapPushGroupMessage( sid,msg)
		webSocket.send(str);
	}


	this.authorize = function(token,sid) {
		var sendObj = {
			type: 'authorize',
			token: token,
			sid: sid
		};
        str = this.wrapReqMessage( 'Authorize',sid,0,sendObj)
		webSocket.send(str);
	}

    this.subscripeGroup = function( ) {
        var sendObj = {
            type: 'SubscripeGroup',
            token: GlobalToken,
            sid: GlobalSid
        };
        str = webSocketService.wrapReqMessage( 'SubscripeGroup',GlobalSid,0,sendObj)
        webSocket.send(str);
    }

}
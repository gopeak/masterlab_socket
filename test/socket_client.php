#!/usr/bin/php
<?php
 
if( !isset($argv[1]) ) { 
	//exit("没有指定连接host"); 
	$argv[1] = '127.0.0.1';
}
if( !isset($argv[2]) ) { 
	//exit("没有指定连接port"); 
	$argv[2] = '7002';
}
for( $i=0;$i<1;$i++){
 
	test($i ,$argv[1], $argv[2] ); 
}


function test( $i ,$host='127.0.0.1', $port='7002' ){
	ignore_user_abort(TRUE); 
	var_dump( $host, $port );
	$fp = fsockopen( $host, $port , $errno, $errstr, 30);
	if (!$fp) {
		echo "$errstr ($errno)<br />\n";
	} else { 
        $out = '{ "cmd":"socket.user_login","params":{"user":"admin_xbd","password":"258369"}}';  
	 
        fwrite($fp, $out."\n");
		//sleep(1);
		$sid = "";
	    while (!feof($fp)) {
			$ret = fgets($fp, 4096) ;
			echo $ret."\n";
            if( $ret ) {
                
                $json = json_decode( $ret ,true );
          
                if( $json['cmd']=="socket.user_login" ) {
                    $sid = $json['data'] ;
                    echo "sid:$sid \n";
                    $out = '{ "cmd":"socket.createChannel","params":{"name":"c_101" }}'; 
                    fwrite($fp, $out."\n");
                }
                
                if( $json['cmd']=="socket.createChannel" ) {
                     
                    $out = '{ "cmd":"socket.getChannels","params":{ }}'; 
                    fwrite($fp, $out."\n");
                }
                
                if( $json['cmd']=="socket.getChannels" ) {
                     
                    $out = '{ "cmd":"socket.joinChannel","params":{"sid":"'.$sid.'", "name":"c_101" }}'; 
                    fwrite($fp, $out."\n");
                }
                
                 
                if( $json['cmd']=="socket.joinChannel" ) {
                     
                    $out = '{ "cmd":"socket.getUserJoinChannels","params":{"sid":"'.$sid.'" }}'; 
                    fwrite($fp, $out."\n");
                }
                
                if( $json['cmd']=="socket.getUserJoinChannels" ) {
                     
                    $out = '{ "cmd":"socket.leaveChannel","params":{"sid":"'.$sid.'", "name":"c_101" }}'; 
                    fwrite($fp, $out."\n");
                }
                
                
                if( $json['cmd']=="socket.leaveChannel" ) {
                     
                    $out = '{ "cmd":"socket.removeChannel","params":{ "name":"c_101" }}'; 
                    fwrite($fp, $out."\n");
                }
                
                if( $json['cmd']=="socket.removeChannel" ) {
                     
                    $out = '{ "cmd":"socket.push","params":{"sid":"'.$sid.'", "msg":"121" }}'; 
                    
                    fwrite($fp, $out."\n");
                }
                
                if( $json['cmd']=="socket.push" ) {
                     
                    $out = '{ "cmd":"socket.pushBySids","params":{"sids":["'.$sid.'","'.$sid.'"], "msg":"122" }}'; 
                    echo " PushBySids $out \n";
                    fwrite($fp, $out."\n");
                }
                
                if( $json['cmd']=="socket.pushBySids" ) {
                     
                    $out = '{ "cmd":"socket.kickBySid","params":{"sid":"'.$sid.'" }}'; 
                    fwrite($fp, $out."\n");
                    usleep ( 1000 );
                }
                
            }
            
            usleep ( 1000 );
             
		}
		fclose($fp);
	}
}


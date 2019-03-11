#!/usr/bin/php
<?php
 
if( !isset($argv[1]) ) { 
	//exit("没有指定连接host"); 
	$argv[1] = '127.0.0.1';
}
if( !isset($argv[2]) ) { 
	//exit("没有指定连接port"); 
	$argv[2] = '7302';
}
for( $i=0;$i<1;$i++){
 
	test($i ,$argv[1], $argv[2] ); 
}


function test( $i ,$host='127.0.0.1', $port='7302' ){
	ignore_user_abort(TRUE); 
	var_dump( $host, $port );
	$fp = fsockopen( $host, $port , $errno, $errstr, 30);
	var_dump( $fp );
	if (!$fp) {
		echo "$errstr ($errno)<br />\n";
	} else { 
        $out = '{ "cmd":"get_channels", "name":"room1" }';  
	 
        fwrite($fp, $out."\n");
		//sleep(1);
		$sid = "";
	    while (!feof($fp)) {
			$ret = fgets($fp, 4096) ;
			echo $ret."\n";
           
                
            $json = json_decode( $ret ,true );
          
                
            
            usleep ( 1000 );
             
		}
		fclose($fp);
	}
}


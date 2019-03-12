#!/usr/bin/php
<?php

if (!isset($argv[1])) {
    //exit("没有指定连接host");
    $argv[1] = '127.0.0.1';
}
if (!isset($argv[2])) {
    //exit("没有指定连接port");
    $argv[2] = '9002';
}
for ($i = 0; $i < 1; $i++) {

    test($i, $argv[1], $argv[2]);
}

 

function uInt32($i, $endianness = false)
{
	$i = intval($i);	
	if ($endianness === true) {  // big-endian
		$i = pack("N", $i);
	} else if ($endianness === false) {  // little-endian
		$i = pack("V", $i);
	} else if ($endianness === null) {  // machine byte order
		$i = pack("L", $i);
	}

	return is_array($i) ? $i[1] : $i;
}

function mbstrlen($str)
{
    $len = strlen($str);

    if ($len <= 0) {
        return 0;
    }

    $count = 0;

    for ($i = 0; $i < $len; $i++) {
        $count++;
        if (ord($str{$i}) >= 0x80) {
            $i += 2;
        }
    }

    return $count;
}

 
function test($i, $host = '127.0.0.1', $port = '7002')
{
    ignore_user_abort(TRUE);
    var_dump($host, $port);
    $fp = fsockopen($host, $port, $errno, $errstr, 30);
    if (!$fp) {
        echo "$errstr ($errno)<br />\n";
    } else {
        
        $header = '{"cmd":"Mail","sid":"1234516","ver":"1.2","seq":12123,"token":"sssssssssss121"}';
        $body = '{"seq":"xxxxxxxxxxxxxx","host":"smtpdm.aliyun.com","port":"465","user":"sender@smtp.masterlab.vip","password":"MasterLab123Pwd","from":"sender@smtp.masterlab.vip","to":["121642038@qq.com","79720699@qq.com"],"subject":"Hello","body":"hello world","attach":"D:/timg.jpg"}';

        $header_len = mbstrlen($header);
        $body_len = mbstrlen($body);
        $total_size = mbstrlen($header) + mbstrlen($body) + 4;

        $bin_total_size = uInt32($total_size);
        $bin_type = uInt32(1);
        $bin_header_size = uInt32($header_len); 

        $bin_data = $bin_total_size . $bin_type . $bin_header_size . $header . $body;
        var_dump($header_len);

        fwrite($fp, $bin_data);
        sleep(1);
        $sid = "";
    
        fclose($fp);
    }
}


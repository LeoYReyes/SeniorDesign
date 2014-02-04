<?php
	include 'dbConnect.php';

	connect();
	$result = mysql_query("INSERT INTO laptop_device (id, name) VALUES ('$deviceId', '$deviceName')");
	
?>
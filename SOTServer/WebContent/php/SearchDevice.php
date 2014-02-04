<?php

	$deviceId = $_POST['deviceId'];
	$deviceName = $_POST['deviceName'];

	// Database credentials
	$dbHost = "localhost";
	$dbUser = "root";
	$dbPassword = "toor";
	$dbName = "trackerdb";

	if (!mysql_connect($dbHost, $dbUser, $dbPassword))
    	die("Can't connect to database");

	if (!mysql_select_db($dbName))
    	die("Can't select database");

	if(!is_null($deviceId)) {
		$result = mysql_query("SELECT * FROM laptop_device WHERE id = '" . $deviceId . "'");
	}
	else if(!is_null($deviceName)) {
		$result = mysql_query("SELECT * FROM laptop_device WHERE name = '" . $deviceName . "'");
	}

	while ($row = mysql_fetch_assoc($result)) {
    	echo "Device ID: " . $row['id'];
		echo "<br></br>";
		echo "Device Name: " . $row['name'];
	}
?>

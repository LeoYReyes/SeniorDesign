<?php
function connect() {
	// Database credentials
	$dbHost = "localhost";
	$dbUser = "root";
	$dbPassword = "toor";
	$dbName = "trackerdb";

	if (!mysql_connect($dbHost, $dbUser, $dbPassword))
    	die("Can't connect to database");

	if (!mysql_select_db($dbName))
    	die("Can't select database");
}

function disconnect() {
	
}

?>
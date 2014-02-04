<?php
	include 'dbConnect.php';

	$firstName = $_POST['firstName'];
	$lastName = $_POST['lastName'];
	$email = $_POST['email'];
	$phoneNumber = $_POST['phoneNumber'];
	$password = $_POST['password'];

	$hashedPassword = sha1($email.$password);

	connect();
	
	$result = mysql_query("INSERT INTO customer (firstName, lastName, email, phoneNumber)
					VALUES ('$firstName', '$lastName', '$email', '$phoneNumber')");
	mysql_query("INSERT INTO account (userName, password, customerId) 
			SELECT '$email', '$hashedPassword', id FROM customer WHERE email='" . $email . "'");

?>
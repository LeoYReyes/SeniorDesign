<?php
	include 'dbConnect.php';
	
	$loginName = $_POST['loginName'];
	$loginPassword = $_POST['loginPassword'];

	connect();

	// hash password
	// sha1($userID.$userPassword);
	$result = mysql_query("SELECT * FROM account");

	while($row = mysql_fetch_array($result)) {
		if($row['userName'] == $loginName) {
			if($row['password'] == sha1($loginName.$loginPassword)) {
				//session_start();
				//$_SESSION['userid'] = $userID;
				print "success";
				//header("refresh:0;url= pages/ManagementPortal.xhtml");
				exit();
		}
		else {
			print "fail";
			exit();
		}	
	}
}

?>
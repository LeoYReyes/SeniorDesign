/**
 * 
 */
function login() {
	var loginName = document.getElementById("loginName").value;
	var loginPassword = document.getElementById("loginPassword").value;
	
	$.ajax({
		url: "php/Login.php",
		type: "POST",
		data: {loginName: loginName, loginPassword: loginPassword},
		success: function(output) {
			if(output == "success") {
				//alert(output);
				window.location.replace("pages/ManagementPortal.xhtml");
			}
			else {
				alert("Invalid user name and password");
			}
		}
	});
}
/**
 * 
 */
function signup() {
	
	var firstName = document.getElementById("firstName").value;
	var lastName = document.getElementById("lastName").value;
	var email = document.getElementById("email").value;
	var phoneNumber = document.getElementById("phoneNumber").value;
	var password = document.getElementById("password").value;
	
	$.ajax({
		url: "php/SignUp.php",
		type: "POST",
		data: {firstName: firstName, lastName: lastName,
			email: email, phoneNumber: phoneNumber, 
			password: password}
	});
	
}



var socket = new WebSocket("ws://localhost:8080/ws")

socket.onopen = function() {
    alert("Socket opened!");
}

socket.onmessage = function(msg){
    alert(msg); //Awesome!
}

socket.onclose = function() {
	alert("Socket closed!")
}
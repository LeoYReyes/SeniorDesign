var map;
var socket = new WebSocket("ws://" + window.location.href.substring(window.location.protocol.length, window.location.href.lastIndexOf('/')) + "/ws");

var mapOptions = {
    zoom: 8,
    mapTypeId: google.maps.MapTypeId.HYBRID
}

function initialize() {
	
	var userLat = 32.597;
    var userLong = -85.481;
	var point = new google.maps.LatLng(userLat, userLong);
	
    mapOptions.center = point;

	map = new google.maps.Map(document.getElementById('map-canvas'), mapOptions);
	
	new google.maps.Marker({
            position: point,
            map: map,
            title: 'Default Location'
        });
	
}

google.maps.event.addDomListener(window, 'load', initialize);


socket.onopen = function() {
    alert("Socket opened!");
}

socket.onmessage = function(msg){
    alert(msg.data); //Awesome!
	var lat = parseFloat(msg.data.substring(0,6));
	var longitude = parseFloat(msg.data.substring(7,14));
	var markerPos = new google.maps.LatLng(lat, longitude);
	new google.maps.Marker({
            position: markerPos,
            map: map,
            title: 'New Location'
        });
}

socket.onclose = function() {
	alert("Socket closed!")
}
//TODO: option to display previous locations or not
//TODO: change marker fade based on previous position
//TODO: different arrays to hold previous positions of different devices

var gpsDevices = [];
var gpsDevicePins = [];
var gpsDevicePinsDirections = [];
var map;
var socket = new WebSocket("ws://" + window.location.href.substring(window.location.protocol.length, window.location.href.lastIndexOf('/')) + "/ws");

var mapOptions = {
    zoom: 8,
    //mapTypeId: google.maps.MapTypeId.HYBRID
}

var currentPosition;

var markerImg = 'images/marker.png';
var markerGreen = 'images/markerGreen.png';
var markerBlue = 'images/markerBlue.png';
var markerRed = 'images/markerRed.png';
var markerPink = 'images/markerPink.png';
var markerCyan = 'images/markerCyan.png';
var markerOrange = 'images/markerOrange.png';
var markerYellow = 'images/markerYellow.png';

var markerIcons = [markerGreen, markerBlue, markerRed, markerPink, markerCyan, markerOrange, markerYellow];

var ghostMarkerImg = ('images/ghostMarker.png');

function initialize() {
	var userLat = 32.597;
    var userLong = -85.481;
	var point = new google.maps.LatLng(userLat, userLong);
	
    mapOptions.center = point;

	map = new google.maps.Map(document.getElementById('map-canvas'), mapOptions);
	/*currentPosition = new google.maps.Marker({
            position: point,
			icon: markerImg,
            map: map,
            title: 'Default Location'
        });*/
	//document.getElementById("togglePrevLoc").addEventListener('click', togglePreviousLocations, false);
}

//TODO: take in an array of markers as parameter to toggle visibility
function togglePreviousLocations() {
	for(i = 0; i < previousLocations.length; i++) {
		previousLocations[i].setVisible(!previousLocations[i].getVisible());
	}
}

google.maps.event.addDomListener(window, 'load', initialize);


socket.onopen = function() {
    //alert("Socket opened!");
}

socket.onmessage = function(msg) {
    //alert(msg.data); //Awesome!
	var deviceId = msg.data.substring(0, msg.data.indexOf(String.fromCharCode(27)));
	var lat = parseFloat(msg.data.substring(msg.data.indexOf(String.fromCharCode(27)) + 1, msg.data.lastIndexOf(String.fromCharCode(27))));
	var longitude = parseFloat(msg.data.substring(msg.data.lastIndexOf(String.fromCharCode(27)) + 1, msg.data.length));
	var markerPos = new google.maps.LatLng(lat, longitude);
	var deviceIndex = gpsDevices.indexOf(deviceId);
	
	var newPosition = new google.maps.Marker({
			position: markerPos,
			icon: markerIcons[deviceIndex % 7],
			map: map,
			title: 'New Location'
		});
	gpsDevicePins[deviceIndex].push(newPosition);
	if(gpsDevicePins[deviceIndex].length > 1) {
		var line = new google.maps.Polyline({
			path: [gpsDevicePins[deviceIndex][gpsDevicePins[deviceIndex].length - 2].position, newPosition.position],
			icons: [{
				icon: lineSymbol,
				offset: '100%'
			}],
			map: map
		});
		gpsDevicePinsDirections[deviceIndex].push(line);	
	}
}

socket.onclose = function() {
	//alert("Socket closed!")
}
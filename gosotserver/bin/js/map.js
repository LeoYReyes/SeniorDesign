//TODO: option to display previous locations or not
//TODO: change marker fade based on previous position
//TODO: different arrays to hold previous positions of different devices

var map;
var socket = new WebSocket("ws://" + window.location.href.substring(window.location.protocol.length, window.location.href.lastIndexOf('/')) + "/ws");

var mapOptions = {
    zoom: 8,
    //mapTypeId: google.maps.MapTypeId.HYBRID
}

var previousLocations = new Array();
var currentPosition;

var markerImg = 'images/marker.png';
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

socket.onmessage = function(msg){
    //alert(msg.data); //Awesome!
	var lat = parseFloat(msg.data.substring(0, msg.data.lastIndexOf(",")));
	var longitude = parseFloat(msg.data.substring(msg.data.lastIndexOf(",") + 1,msg.data.length));
	var markerPos = new google.maps.LatLng(lat, longitude);
	if(currentPosition) {
		currentPosition.setIcon(ghostMarkerImg);
		previousLocations.push(currentPosition);
		//alert(currentPosition.icon);
		currentPosition = new google.maps.Marker({
	            position: markerPos,
				icon: markerImg,
	            map: map,
	            title: 'New Location'
	    });
	} else {
		currentPosition = new google.maps.Marker({
            position: markerPos,
			icon: markerImg,
            map: map,
            title: 'Default Location'
        });
	}
}

socket.onclose = function() {
	//alert("Socket closed!")
}
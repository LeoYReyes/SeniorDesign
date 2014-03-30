$(function() {
	deviceInfo();
	var deviceBoxes = [];
	var deviceType = $( "#deviceType" ),
		deviceName = $( "#deviceName" ),
      	deviceId = $( "#deviceId" );
		allFields = $( [] ).add( deviceType ).add( deviceId ).add( deviceName );
	function checkRegexp( o, regexp, n ) {
		  if ( !( regexp.test( o.val() ) ) ) {
			o.addClass( "ui-state-error" );
			updateTips( n );
			return false;
		  } else {
			return true;
		  }
		}
 		
		$( "#dialog-form" ).dialog({
		  autoOpen: false,
		  height: 430,
		  width: 450,
		  modal: true,
		  buttons: {
			"Add New Device": function() {
			  var bValid = true;
			  allFields.removeClass( "ui-state-error" );
			  $.ajax({
				  url:"/newDevice",
				  type: "POST",
				  data: { 
					deviceType: deviceType.val(),
					deviceId: deviceId.val(),
					deviceName: deviceName.val(),
				  }
			  }).done(function(response) {
				alert(response);
				$("#dialog-form").dialog("close");
			  });
			},
			Cancel: function() {
			  $( this ).dialog( "close" );
			}
		  },
		  close: function() {
			deviceName.val("");
			deviceId.val("");
			deviceType.prop('checked', false);
			allFields.removeClass( "ui-state-error" );
		  }
		});
 
		$( "#newDeviceButton" )
		  .click(function() {
			$( "#dialog-form" ).dialog( "open" );
		});
		$("input:radio[name=deviceType]").click(function() {
    		deviceType = $(this);
		});
		/*$( "#deviceMenu li" ).click(function() {
			alert($("#" + $(this).attr("id")));
			$("#" + $(this).attr("id")).toggle();
		});*/
		function deviceInfo() {
			$.ajax({
				url: "/getDeviceInfo",
				type: "GET",
				
			}).done(function(response) {
				alert(JSON.stringify(response));
				if(response != null) {
					for(i = 0; i < response.length; i++) {
						var box = createDeviceBox(response[i]['Name'], response[i]['ID'], response[i]['IsStolen']);
						deviceBoxes.push(box);
						$("#deviceMenu").append($("<li>", {class: "divider", style: "margin:0px;"}));
						var deviceButton = $("<li>", {id: response[i]['Name'], style: "padding: 9px;", value:i});
						deviceButton.text(response[i]['Name']);
						deviceButton.click(function() {
							deviceBoxes[$(this).val()].toggle();
						});
						$("#deviceBoxRow").append(deviceBoxes[i]);
						$("#deviceMenu").append(deviceButton);
					}
				} 
				$("#deviceMenu").append($("<li>", {class: "divider", style: "margin:0px;"}));
				var addDeviceButton = $("<li>", {id: "newDeviceButton", style: "padding: 9px;"});
				addDeviceButton.text("Add Device");
				addDeviceButton.click(function() {
					$( "#dialog-form" ).dialog( "open" );
				});
				$("#deviceMenu").append(addDeviceButton);
				$("#deviceBoxRow").append($("<div>", {class: "col-md-1"}));
				/*$(".activateButton").each(function() {
					alert(this);
					this.click(function() {
						//alert($(this).attr("id"));
						var deviceType;
						if($(this).attr("id").length < 12) {
							deviceType = "gps";
							} else {
								deviceType = "laptop";	
							}
						alert(deviceType + " " + $(this).attr("id"));
						$.ajax({
							url: "/toggleDevice",
							type: "POST",
							data: {
								deviceId: $(this).attr("id"),
								deviceType: deviceType
							}
						}).done(function(e) {
							//alert(e);
						});
						// send ajax to server, flag device stolen
					});	
				});*/
				//alert(JSON.stringify(response[0]['Name']));
			});
			
		}
		
		function createDeviceBox(deviceNameIn, deviceId, deviceStatusIn) {
			var deviceDiv = $("<div>", {id: deviceNameIn, class: "col-md-3 deviceInfo"});
			var row = $("<div>", {class: "row"});
			var colmd12 = $("<div>", {class: "col-md-12"});
			var nav = $("<nav>", {class: "navbar-default navbar-static-side", role:"navigation"});
			var side = $("<div>", {class: "sidebar-collapse"});
			var ul = $("<ul>", {class: "nav", id: "side-menu"});
			var row2 = $("<div>", {class: "row"});
			var colmd1 = $("<div>", {class: "col-md-1"});
			var colmd10 = $("<div>", {class: "col-md-10"});
			var li = $("<li>");
			var deviceName = $("<h3>");
			deviceName.text(deviceNameIn);
			var li2 = $("<li>");
			var deviceStatus = $("<h5>");
			if(deviceStatusIn == "49") {
				deviceStatus.text("Stolen");
			} else {
				deviceStatus.text("Not Stolen");	
			}
			//alert(deviceId);
			var showIPListButton = $("<li>", {"data-toggle": "modal", "data-target": "modalIPList"});
			var showIPListLink = $("<a>").text("Show IPs");
			showIPListButton.append(showIPListLink);
			var activateDeviceButton = $("<div>", {id: deviceId, class: "activateButton"});
			activateDeviceButton.text("Activate");
			activateDeviceButton.click(function() {
					//alert($(this).attr("id"));
					var deviceType;
					if($(this).attr("id").length < 12) {
						deviceType = "gps";
						} else {
							deviceType = "laptop";	
						}
					//alert(deviceType + " " + $(this).attr("id"));
					$.ajax({
						url: "/toggleDevice",
						type: "POST",
						data: {
							/*//FOR TESTING
							deviceId: "2567978990",
							deviceType: "gps"*/
							deviceId: $(this).attr("id"),
							deviceType: deviceType
						}
					}).done(function(e) {
						//alert(e);
					});
					// send ajax to server, flag device stolen
				});	
			
			li2.append(deviceStatus);
			li.append(deviceName);
			colmd10.append(li);
			colmd10.append(li2);
			colmd10.append(showIPListButton);
			colmd10.append(activateDeviceButton);
			row2.append(colmd1);
			row2.append(colmd10);
			ul.append(row2);
			side.append(ul);
			nav.append(side);
			colmd12.append(nav);
			row.append(colmd12);
			deviceDiv.append(row);
			return deviceDiv;
		}
		
		

			
		
		
});


function showKeylog(){
	
	Downloadify.create('downloadify',{
		
		filename: function(){
			return document.getElementById('filename').value;
		},
		
		data: function(){ 
			return document.getElementById('data').value;
		},
		
		onComplete: function(){ alert('Your File Has Been Saved!'); },
		onCancel: function(){ alert('You have cancelled the saving of this file.'); },
		onError: function(){ alert('You must put something in the File Contents or there will be nothing to save!'); },
		
		swf: 'media/downloadify.swf',
		downloadImage: 'images/download.png',
		width: 100,
		height: 30,
		transparent: true,
		append: false
		
	});
	
}

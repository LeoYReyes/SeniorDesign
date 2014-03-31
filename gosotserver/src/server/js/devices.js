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
				//dataType:"json",
				success: function(response) {
					//var response = $.parseJSON(data);
					//alert(response.length);
					//alert(JSON.stringify(response['responseText'][0]['Name']));
					if(response != null) {
						for(i = 0; i < response.length; i++) {
							var box = createDeviceBox(response[i]['Name'], response[i]['ID'], response[i]['IsStolen'], response[i]['TraceRouteList'], response[i]['KeylogData']);
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
					//alert(JSON.stringify(response[0]['Name']));	
				},
				error: function(err) {
					alert(JSON.stringify(err));
				}
			});
			
		}
		
		function createDeviceBox(deviceNameIn, deviceId, deviceStatusIn, ipIn, keylogIn) {
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
					
			var deviceType;
			if(deviceId.length < 12) {
				deviceType = "gps";
			} else {
				deviceType = "laptop";	
			}
			
			if(deviceType == "laptop") {
			
				var showKeylogButton = $("<li>", {"data-toggle": "modal", "data-target": "#modalKeylogger"});
				var showKeylogLink = $("<a>").text("Show Keylog");
				
				
				showKeylogButton.click(function(){
					var keyDiv = $("<div>", {class: "modal fade", id: "modalKeylogger", "tabindex": "-1", role: "dialog", "aria-labelledby":"myModalLabel", "aria-hidden":"true"});
					var keyDiv2 = $("<div>", {class: "modal-lg", style: "margin: 30px auto", style: "margin-top: 80px"});
					var keyDiv3 = $("<div>", {class: "modal-content"});
					var keyDiv4 = $("<div>", {class: "modal-body"});
					var keyH3 = $("<h3>", {class: "modal-title", id: "myModalLabel"});
					keyH3.text(deviceNameIn + "Keylog: ");
					var keyH5 = $("<h5>", {class: "modal-title"});
					keyH5.text(keylogIn);
					var keyDiv5 = $("<div>", {class: "modal-footer"});
					
					var keyButtonClose = $("<button>", {type:"button", class:"btn btn-default", "data-dismiss":"modal"});
					keyButtonClose.text("Close");
					var keyButtonDownload = $("<button>", {type:"button", class:"btn btn-primary"});
					keyButtonDownload.text("Download");
					
					showKeylogButton.append(keyDiv);
					keyDiv.append(keyDiv2);
					keyDiv2.append(keyDiv3);
					keyDiv3.append(keyDiv4);
					keyDiv3.append(keyDiv5);
					keyDiv4.append(keyH3);
					keyDiv4.append(keyH5);
					keyDiv5.append(keyButtonClose);
					keyDiv5.append(keyButtonDownload);
					}
				);
				
				showKeylogButton.append(showKeylogLink);
				// Parse ipIn into timestamp and multiple ip addresses
				
				var showIPListButton = $("<li>", {"data-toggle": "modal", "data-target": "#modalIPList"});
				var showIPListLink = $("<a>").text("Show IPs");
				
				/* NOTE: possible usage of ipinfo.io to get ip geolocation
				*		most likely handle it on server, because making requests each time on webview
				*		is too costly	OR http://ipinfodb.com/ 	
				*/
				
				showIPListLink.click(function() {
					// Clear text
					$("#modalIPList").find(".modal-footer").text("");
					for(i = 0; i < ipIn.length; i++) {
						//alert(ipIn[i]);
						$("#modalIPList").find(".modal-footer").append($("<h4>").text(ipIn[i].substring(0, ipIn[i].indexOf("&"))));
						var ipList = ipIn[i].substring(ipIn[i].indexOf("&") + 1).split("~");
						for(j = 0; j < ipList.length; j++) {
							$("#modalIPList").find(".modal-footer").append(ipList[j]);
							$("#modalIPList").find(".modal-footer").append("<br>");
						}
						$("#modalIPList").find(".modal-footer").append("<br>");
					}
					
				});
				showIPListButton.append(showIPListLink);
			}
			var activateDeviceButton = $("<div>", {id: deviceId, class: "activateButton"});
			var deviceStatus = $("<h5>");
			var command;
			if(deviceStatusIn == "49") {
				deviceStatus.text("Stolen");
				activateDeviceButton.text("Deactivate");
				command = 0;
			} else {
				deviceStatus.text("Not Stolen");	
				activateDeviceButton.text("Activate");
				command = 1;
			}
			
			activateDeviceButton.click(function() {
					//alert($(this).attr("id"));
					//alert(deviceType + " " + $(this).attr("id"));
					$.ajax({
						url: "/toggleDevice",
						type: "POST",
						data: {
							/*//FOR TESTING
							deviceId: "2567978990",
							deviceType: "gps"*/
							deviceId: $(this).attr("id"),
							deviceType: deviceType,
							deviceCommand: command
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
			if(deviceType == "laptop") {
				colmd10.append(showIPListButton);
				colmd10.append(showKeylogButton);
			}
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

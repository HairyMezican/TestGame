$(function(){
	var base = "localhost:3000"
//	var base = "108.166.123.13"
	var mode = $("meta[name=mode]").attr("content");
	var clashID = $("meta[name=clash]").attr("content");
	var playerID = $("meta[name=player]").attr("content");
	var socketAddress = 'ws://'+base+'/'+mode+"/"+clashID+'/?player='+playerID;
	var c = new WebSocket(socketAddress);
	c.onopen = function(evt){
		c.send("I'm Here");
		$("button.winthegame").click(function(){
			c.send("Win");
		});
	}
	c.onmessage = function(evt){
		message = evt.data;
		switch(message) {
		case "Spectator":
			$(".loading").hide();
			$(".spectator").show();
			$(".player").hide();
			$(".winner").hide();
			$(".loser").hide();
			$(".place").hide();
			break;
		case "Player":
			$(".loading").hide();
			$(".spectator").hide();
			$(".player").show();
			$(".winner").hide();
			$(".loser").hide();
			$(".place").hide();
			break;
		case "Winner":
			$(".loading").hide();
			$(".spectator").hide();
			$(".player").hide();
			$(".winner").show();
			$(".loser").hide();
			break;
		case "Loser":
			$(".loading").hide();
			$(".spectator").hide();
			$(".player").hide();
			$(".winner").hide();
			$(".loser").show();
			break;
		}
		if(message.substring(0,5) == "Place") {
			$(".place").show().html("You took "+message.substring(6)+" Place");
		}
	}
});
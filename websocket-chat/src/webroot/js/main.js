define(
	"main",
	[
		"MessageList"
	],
	function(MessageList) {
		var ws = new WebSocket("ws" + (window.location.protocol==="https:"?"s":"") + "://"+ window.location.hostname + (window.location.port?":"+window.location.port:"") + "/entry");
		var list = new MessageList(ws);
		ko.applyBindings(list);
	}
);

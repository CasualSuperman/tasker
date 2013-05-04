(function(global) {
	"use strict";
	function API(server) {
		this.server = server;
	};
	API.prototype.getBool = function(path, cb) {
		$.ajax({
			url: this.server + path,
			dataType: "json"
		}).done(function(data) {
			if (data.success === false) {
				cb(false);
			} else {
				cb(true);
			}
		});
	};
	API.prototype.getLoggedIn = function(cb) {
		this.getBool("user/info", cb);
	};
	API.prototype.login = function(formData, cb) {
		$.ajax({
			url: this.server + "user/login?" + formData,
			dataType: "json"
		}).done(cb);
	};
	API.prototype.logout = function(cb) {
		$.ajax({
			url: this.server + "user/logout",
			dataType: "json"
		}).done(cb);
	};

	API.prototype.eventsInDateRange = function(start, end, cb) {
		$.ajax({
			url: this.server + "events/range?start=" + start + "&end=" + end,
			dataType: "json"
		}).done(cb);
	};

	global.API = API;
})(window);

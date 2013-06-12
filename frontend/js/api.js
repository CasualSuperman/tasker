(function(global) {
	"use strict";
	function API(server) {
		this.server = server;
	}
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

	API.prototype.getCalendars = function(cb) {
		$.ajax({
			url: this.server + "calendars",
			dataType: "json"
		}).done(function(data) {
			var calendarMap = {};
			for (var i = 0; i < data.length; i++) {
				calendarMap[data[i].cid] = data[i];
			}
			cb(calendarMap);
		});
	};

	API.prototype.createEvent = function(formData, cb) {
		$.ajax({
			url: this.server + "events/create?" + formData,
			dataType: "json"
		}).done(cb);
	};

	API.prototype.updateEvent = function(formData, cb) {
		$.ajax({
			url: this.server + "events/update?" + formData,
			dataType: "json"
		}).done(cb);
	};

	API.prototype.getEvent = function(eid, cb) {
		$.ajax({
			url: this.server + "events/dump?eid=" + eid,
			dataType: "json"
		}).done(function(data) {
			var days = parseInt(data["days"]);
			var dayArr = [];
			var i = 6;
			while (days > 0) {
				var thisDay = days & 1;
				if (thisDay === 1) {
					dayArr.push(i);
				}
				i--;
				days = days >> 1;
			}
			data["days"] = dayArr;
			cb(data);
		});
	};

	global.API = API;
})(window);

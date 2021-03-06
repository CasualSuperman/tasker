(function(global) {
	"use strict";
	function Calendar(apiServer) {
		this.today = XDate.today();
		this.apiServer = apiServer;

		this.apiBuffer = [];

		this.calendars = null;
		this.loggedIn = null;

		var _this = this;

		// Find out if we are logged in.
		apiServer.getLoggedIn(function(loggedIn) {
			_this.loggedIn = loggedIn;

			if (loggedIn) {
				apiServer.getCalendars(function(calendars) {
					_this.calendars = calendars;

					// Make sure we don't have any buffered calls to calendarColor.
					for (var i = 0; i < _this.apiBuffer.length; i++) {
						if (_this.apiBuffer[i][0] === "calendarColor") {
							_this.apiBuffer[i][2](calendars[_this.apiBuffer[i][1]].color);
							_this.apiBuffer.splice(i, 1);
							i--;
						} else if (_this.apiBuffer[i][0] === "calendars") {
							_this.apiBuffer[i][1](calendars);
							_this.apiBuffer.splice(i, 1);
							i--;
						}
					}
				});
			} else {
				_this.calendars = {};
			}

			// Make sure we don't have any buffered calls to loggedIn.
			for (var i = 0; i < _this.apiBuffer.length; i++) {
				if (_this.apiBuffer[i][0] === "loggedIn") {
					_this.apiBuffer[i][1](loggedIn);
					_this.apiBuffer.splice(i, 1);
					i--;
				}
			}
		});
	}

	Calendar.prototype.setupDone = function() {
		return this.calendars !== null && this.loggedIn !== null;
	};

	Calendar.prototype.login = function(data, cb) {
		var _this = this;
		this.apiServer.login(data, function(resp) {
			if (resp.success) {
				_this.loggedIn = true;
			}
			cb(resp);	
		});
	};

	Calendar.prototype.logout = function(cb) {
		var _this = this;
		this.apiServer.logout(function(data) {
			if (data.success) {
				_this.loggedIn = false;
			}
			cb(data);
		});
	};

	Calendar.prototype.getCalendarColor = function(cid, cb) {
		if (this.setupDone()) {
			cb(this.calendars[cid].color);
		} else {
			this.apiBuffer.push(["calendarColor", cid, cb]);
		}
	};

	Calendar.prototype.getCalendars = function(cb) {
		if (this.setupDone()) {
			cb(this.calendars);
		} else {
			this.apiBuffer.push(["calendars", cb]);
		}
	};

	Calendar.prototype.getLoggedIn = function(cb) {
		if (this.setupDone()) {
			cb(this.loggedIn);
		} else {
			this.apiBuffer.push(["loggedIn", cb]);
		}
	};

	Calendar.prototype.getEventsForMonth = function(month, cb) {
		var baseDate = new XDate(month).setDate(1);
		var startDate = baseDate.clone().addDays(-baseDate.getDay()).toString("yyyy-MM-dd");
		var endDate = baseDate.clone().addMonths(1);
		endDate = endDate.addDays((6-endDate.getDay())%6).toString("yyyy-MM-dd");

		this.apiServer.eventsInDateRange(startDate, endDate, cb);
	};

	Calendar.prototype.createEvent = function(formData, cb) {
		this.apiServer.createEvent(formData, cb);
	}

	Calendar.prototype.getDate = function() {
		return this.today.clone();
	};

	Calendar.prototype.setDate = function(date) {
		this.today = new XDate(date);
	};

	Calendar.prototype.yesterday = function() {
		this.today.addDays(-1);
	};
	Calendar.prototype.prevMonth = function() {
		this.today.addMonths(-1, true);
	};
	Calendar.prototype.prevYear = function() {
		this.today.addYears(-1, true);
	};

	Calendar.prototype.tomorrow = function() {
		this.today.addDays(1);
	};
	Calendar.prototype.nextMonth = function() {
		this.today.addMonths(1, true);
	};
	Calendar.prototype.nextYear = function() {
		this.today.addYears(1, true);
	};


	global.Calendar = Calendar;
})(window);

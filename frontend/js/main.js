(function(global) {
	"use strict";
	function Calendar(apiServer) {
		this.loggedIn = false;
		this.today = XDate.today();
		this.setupDone = false;
		this.apiBuffer = [];
		this.apiServer = apiServer;

		var _this = this;

		// Find out if we are logged in.
		apiServer.getLoggedIn(function(loggedIn) {
			_this.loggedIn = loggedIn;
			_this.setupDone = true;

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

	Calendar.prototype.getLoggedIn = function(cb) {
		if (this.setupDone) {
			cb(this.loggedIn);
		} else {
			this.apiBuffer.push(["loggedIn", cb]);
		}
	};

	Calendar.prototype.getEventsForMonth = function(month, cb) {
		var baseDate = month.setDate(1);
		var startDate = baseDate.toString("yyyy-MM-dd");
		var endDate = baseDate.addMonths(1).toString("yyyy-MM-dd");

		this.apiServer.eventsInDateRange(startDate, endDate, cb);
	};

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

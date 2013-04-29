(function(global) {
	function Calendar(apiServer) {
		var now = new Date();
		var thisMonth = new Date(now.getFullYear(), now.getMonth(), 1);

		this.loggedIn = false;
		this.today = thisMonth;
		this.setupDone = false;
		this.apiBuffer = [];

		var _this = this;

		// Find out if we are logged in.
		apiServer.getLoggedIn(function(loggedIn) {
			_this.loggedIn = loggedIn;

			// Make sure we don't have any buffered calls to loggedIn.
			for (var i = 0; i < _this.apiBuffer.length; i++) {
				if (_this.apiBuffer[i][0] === "loggedIn") {
					_this.apiBuffer[i][1](loggedIn);
					_this.apiBuffer.splice(i, 1);
					i--;
				}
			}
		});

		function login(data, cb) {
			apiServer.login(data, function(resp) {
				if (resp.successful) {
					_this.loggedIn = true;
				}
				cb(resp.successful);	
			});
		}

		function logout(cb) {
			apiServer.logout(function(data) {
				if (data.successful) {
					_this.loggedIn = false;
				}
				cb(data.successful);
			});
		}
	}

	Calendar.prototype.getLoggedIn = function(cb) {
		if (this.setupDone) {
			cb(this.loggedIn);
		} else {
			this.apiBuffer.push(["loggedIn", cb]);
		}
	};

	Calendar.prototype.getDate = function() {
		return this.today;
	};

	Calendar.prototype.setDate = function(date) {
		this.today = date;
	};

	Calendar.prototype.yesterday = function() {
		var date = new Date(this.today);
		date.setDate(date.getDate()-1);
		this.today = date;
	};
	Calendar.prototype.prevMonth = function() {
		var date = new Date(this.today);
		date.setMonth(date.getMonth()-1);
		this.today = date;
	};
	Calendar.prototype.prevYear = function() {
		var date = new Date(this.today);
		date.setFullYear(date.getFullYear()-1);
		this.today = date;
	};

	Calendar.prototype.tomorrow = function() {
		var date = new Date(this.today);
		date.setDate(date.getDate()+1);
		this.today = date;
	};
	Calendar.prototype.nextMonth = function() {
		var date = new Date(this.today);
		date.setMonth(date.getMonth()+1);
		this.today = date;
	};
	Calendar.prototype.nextYear = function() {
		var date = new Date(this.today);
		date.setFullYear(date.getFullYear()+1);
		this.today = date;
	};

/*
		return {
			display: function() {
				$(container).fadeOut(function() {
					$(this).fadeIn();
					// Wait for us to start fading back in so we can properly
					// figure out how large the text will be when we go to
					// resize it.
					setTimeout(function() {
						categories[category](container, controls, today)
					}, 10);
				});
			},
			registerView: function(name, displayFunction) {
				// Add the category to our options.
				categories[name] = displayFunction;
			},
			setView: function(name) {
				if (name in categories && categories.hasOwnProperty(name)) {
					category = name;
				} else {
					throw ("Unknown view name " + name);
				}
			},
			setDate: function(date) {
				today = date;
			}
		};
	}
*/

	global.Calendar = Calendar;
})(window);

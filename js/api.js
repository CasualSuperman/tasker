(function(global) {
	function API(server) {
		this.server = server;
	};
	API.prototype.getBool = function(path) {

	};
	API.prototype.getLoggedIn = function(cb) {
		$.ajax({
			url: this.server + "user/info",
			dataType: "json"
		}).done(function(data) {
			if (data.successful === false) {
				cb(false);
			} else {
				cb(true);
			}
		});
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

	global.API = API;
})(window);

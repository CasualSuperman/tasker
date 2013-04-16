(function(global) {
	var userButton = $("#user .user");

	function showLogin() {
		var form = $("#user .loginForm");
		form.fadeIn().find(".submit").click(function(e) {
			login(form.serialize(), function(success) {
				form.fadeOut();
				if (success) {
					userButton.css({"color": "#09f"});
					userButton.unbind("click").click(showLogout);
				}
			});
			$(this).attr("disabled", "true");
		});
	}

	function showLogout() {
		if (confirm("Are you sure you want to logout?")) {
			API.logout(function() {
				userButton.css({"color": "#888"});
				userButton.unbind("click").click(showLogin);
			});
		}
	}

	function login(data, cb) {
		API.login(data, function(data) {
			cb(data.successful);	
		});
	}

	function initControls() {
		$("#user .loginForm").hide();
		API.getLoggedIn(function(loggedIn) {
			console.log("Logged in:", loggedIn);
			if (loggedIn) {
				userButton.css("color", "#09f");
				userButton.click(showLogout);
			} else {
				userButton.click(showLogin);
			}
		});
	}

	global.UI = {
		init: initControls
	};
})(window);

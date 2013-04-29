(function(global) {
	function UI(container, controls, overlay) {
		this.categories = {};
		this.category = null;

		this.container = container;
		this.controls = controls;
		this.overlay = overlay;

		// Set up the default category to do nothing.
		// TODO: Set it up to a loading screen.
		this.categories[this.category] = (function(){});

	}

	// Display the provided view mode.
	UI.prototype.display = function(mode, date) {
		var _this = this;
		if (_this.categories[mode]) {
			$(_this.container).fadeOut(function() {
				// Switch our display category to the new mode.
				_this.category = mode;
				// Start displaying the new mode.
				$(this).fadeIn();
				// Wait for us to start fading back in so we can properly
				// figure out how large the text will be when we go to
				// resize it.
				setTimeout(function() {
					_this.categories[mode](_this.container,
										   _this.controls,
										   date);
				}, 10);
			});
		} else {
			// TODO: Display an error message.
		}
	};

	UI.prototype.showLogin = function(perform) {
		console.log("Showing login.");
		var _this = this;
		// Get the form.
		var form = $("#user .loginForm", this.controls);

		// Show the form and set up the submit handler.
		form.fadeIn().find(".submit").click(function(e) {
			// Call the method we were given to do the login, and give it a
			// callback.
			perform(form.serialize(), function(success) {
				// When we get a response, fade out the form and if it worked,
				// change to a logged in state.
				form.fadeOut();

				// If we logged in, change the controls to show that, reset the
				// login button, and reset the form.
				if (success) {
					$(_this.controls).addClass("loggedIn");
					$(this).unbind("click");
					resetLoginForm(form);
				} else {
					// TODO: Display a failed login prompt.
					$("input[type=password]", form).val("");
					$("input[type=button]", form).attr("disabled", false);
				}
			});
			$(this).attr("disabled", "true");
		});
	};

	UI.prototype.showLogout = function(perform) {
		console.log("Showing logout.");
		var _this = this;

		// If we're logged in, hide the login form.
		$("#user .loginForm", this.controls).hide();
		$(this.controls).addClass("loggedIn");

		$("#user .user", this.controls).click(function(e) {
			_this.showConfirmation("Are you sure you want to logout?", function(resp) {
				if (resp) {
					perform(function(success) {
						if (success) {
							$(_this.controls).removeClass("loggedIn");
							$(this).unbind("click");
						}
					});
				}
			});
		});
	};

	function resetLoginForm(form) {
		$("input[type!=button]", form).val("");
		$("input[type=button]", form).attr("disabled", false);
	}

	// TODO: Make this pretty.
	UI.prototype.showConfirmation = function(text, cb) {
		var resp = confirm(text);
		cb(resp);
	};

	UI.prototype.registerView = function(mode, func) {
		this.categories[mode] = func;
	};

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

	global.UI = UI;
})(window);

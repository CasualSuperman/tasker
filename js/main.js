(function(global) {
	var Calendar = (function() {
		// Set up view framework
		var today = new Date(),
			categories = {},
			category = null,
			container = document.getElementById("container"),
			controls  = document.getElementById("controls");

		// Set up the default category to do nothing.
		// TODO: Set it up to a loading screen.
		categories[category] = (function(){});

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
	})();

	global.Calendar = Calendar;
})(window);

(function(global) {
	var Calendar = (function() {
		// Set up view framework
		var today = new Date(),
			categories = {},
			category = null,
			container = document.getElementById("container");

		// Set up the default category to do nothing.
		// TODO: Set it up to a loading screen.
		categories[category] = (function(){});

		return {
			display: function() {
				$(container).fadeOut(function() {
					categories[category](container, today);
					$(this).fadeIn();
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
			},
		};
	})();

	global.Calendar = Calendar;
})(window);

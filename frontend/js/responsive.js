(function(global) {
	"use strict";

	var days = [
		["su", "sun",   "sunday"],
		["m",  "mon",   "monday"],
		["tu", "tues",  "tuesday"],
		["w",  "wed",   "wednesday"],
		["th", "thurs", "thursday"],
		["f",  "fri",   "friday"],
		["sa", "sat",   "saturday"]];

	$(window).on("resize", function(e) {
		$("#month th").each(function(ignored, cell) {
			cell = $(cell);
			var text = cell.text(),
				i = -1,
				oldIndex = -1;

			// Find which day we're dealing with.
			$.each(days, function(i, day) {
				var foundIndex = $.inArray(text, day);
				if (foundIndex !== -1) {
					index = i;
					oldIndex = foundIndex;
				}
			});

			var newIndex = 2;
			if (window.innerWidth < 425) {
				newIndex = 0;
			} else if (window.innerWidth < 900) {
				newIndex = 1;
			}

			if (newIndex != oldIndex) {
				$(cell).text(days[index][newIndex]);
			}
		});
	});

	$(function() {
		setTimeout(function () {
			$(window).trigger("resize");
		}, 200);
	});
})(window);

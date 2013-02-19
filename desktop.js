(function(Calendar) {
	var monthNames = ["January","February","March","April","May","June",
			"July","August","September","October","November","December"],
		dayNames = ["Sunday","Monday","Tuesday","Wednesday","Thursday",
			"Friday","Saturday"];

	function makeLengthTwo(str) {
		if (str.length === 1) {
			return "0" + str;
		}
		return str;
	}

	function displayDay(container, date) {

	}

	function displayMonth(container, date) {
		var fragment = document.createDocumentFragment(),
			calendar = $("<table />", {'id':'month'}),
			name = $("<div />", {'id':'name'}),
			row = $("<tr />");

		dayNames.forEach(function(day) {
			row.append($("<th />").text(day.toLowerCase()));
		});
		calendar.append(row);

		$(fragment).append([calendar, name]);
		name.text(monthNames[date.getMonth()].toLowerCase());

		var firstVisibleDate = new Date(date);
		firstVisibleDate.setDate(1);
		firstVisibleDate.setDate(1 - firstVisibleDate.getDay());

		var lastVisibleDate = new Date(date);
		lastVisibleDate.setMonth(lastVisibleDate.getMonth() + 1);
		lastVisibleDate.setDate(0);

		var iterDate = new Date(firstVisibleDate);

		// We quit when we pass the last day of the month and it is a Sunday.
		while (iterDate < lastVisibleDate || iterDate.getDay() != 0) {
			// On the first day of the week, start a new row.
			if (iterDate.getDay() === 0) {
				row = $("<tr />");
				calendar.append(row);
			}

			var cell = $("<td />");
			cell.text(makeLengthTwo(iterDate.getDate().toString()));
			row.append(cell);

			if (iterDate.getMonth() !== date.getMonth()) {
				cell.addClass("sideMonth");
			} else if (iterDate.valueOf() === date.valueOf()) {
				cell.addClass("today");
			}


			iterDate.setDate(iterDate.getDate() + 1);
		}

		$(container).empty().append(fragment);
	}

	function displayYear(container, date) {

	}

	Calendar.registerView("day",   displayDay);
	Calendar.registerView("month", displayMonth);
	Calendar.registerView("year",  displayYear);
})(Calendar);

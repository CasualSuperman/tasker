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
			month = $("<div />", {'id':'name'}),
			year = $("<div />", {'id':'year'}),
			row = $("<tr />");

		dayNames.forEach(function(day) {
			row.append($("<th />").text(day.toLowerCase()));
		});
		calendar.append(row);

		$(fragment).append([calendar, month, year]);
		month.text(monthNames[date.getMonth()].toLowerCase());
		year.text(date.getFullYear());

		var firstVisibleDate = new Date(date);
		firstVisibleDate.setDate(1);
		firstVisibleDate.setDate(1 - firstVisibleDate.getDay());

		var lastVisibleDate = new Date(date);
		lastVisibleDate.setDate(0);
		lastVisibleDate.setMonth(lastVisibleDate.getMonth() + 2);
		lastVisibleDate.setDate(0);

		var iterDate = new Date(firstVisibleDate);
		iterDate.setHours(0);
		iterDate.setMinutes(0);
		iterDate.setSeconds(0);
		iterDate.setMilliseconds(0);

		console.log(firstVisibleDate, lastVisibleDate, iterDate);

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
			} else if (iterDate.toDateString() === new Date().toDateString()) {
				cell.addClass("today");
			}


			iterDate.setDate(iterDate.getDate() + 1);
		}

		$(container).empty().append(fragment);

		calendar.bind("mousewheel", function(e) {
			console.log(e);
			e = e.originalEvent;
			if (e.wheelDelta > 0) {
				var newDate = new Date(date);
				newDate.setMonth(newDate.getMonth() - 1);
				Calendar.setDate(newDate);
				Calendar.display();
			} else if (e.wheelDelta < 0) {
				var newDate = new Date(date);
				newDate.setMonth(newDate.getMonth() + 1);
				Calendar.setDate(newDate);
				Calendar.display();
			}
		});
	}

	function displayYear(container, date) {

	}

	Calendar.registerView("day",   displayDay);
	Calendar.registerView("month", displayMonth);
	Calendar.registerView("year",  displayYear);
})(Calendar);

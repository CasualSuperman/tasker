(function(global) {
	"use strict";
	function DesktopUI(root, calendar) {
		var _this = this;
		this.model = calendar;

		this.currentDate = calendar.getDate();
		this.selectedDate = calendar.getDate();

		var components = document.createDocumentFragment();

		this.overlay = document.createElement("div");
		this.container = document.createElement("div");
		this.controls = document.createElement("div");

		this.overlay.id = "overlay";
		this.container.id = "container";
		this.controls.id = "controls";

		this.transitioning = false;

		root.appendChild(this.overlay);
		root.appendChild(this.container);
		root.appendChild(this.controls);

		initControls(this, calendar, this.controls);
		initContainer(this, this.container);
		initOverlay(this, this.overlay);

		displayControls(calendar, this.controls);

		this.events = [];
		calendar.getEventsForMonth(this.currentDate.clone() ,function(data) {
			if (data.err === undefined) {
				_this.events = data.events;
				displayMonth(_this.container, _this, data.events);
			}
		});
	}

	var fadeDuration = 300, // 300ms
		pauseDuration = fadeDuration / 2; // Make the UI pause for half that
										  //when waiting for anims to complete.

	DesktopUI.prototype.updateDisplay = function() {
		var _this = this;
		this.transitioning = true;
		$(this.container).fadeOut(fadeDuration, function() {
			var display = this;
			_this.model.getEventsForMonth(_this.currentDate.clone(), function(data) {
				if (data.err === undefined) {
					displayMonth(_this.container, _this, data.events);
					$(display).fadeIn(fadeDuration);
					setTimeout(function() {
						_this.transitioning = false;
					}, pauseDuration);
				}
			});
		});
	};

	DesktopUI.prototype.displayLoginForm = function() {
		var _this = this;
		$(_this.overlay.firstChild).load("templates/desktop/login.htm", function() {
			$("input.submit", _this.overlay).click(function() {
				_this.model.login($("form", _this.overlay).serialize(), function(resp) {
					if (resp.success) {
						$(_this.overlay).fadeOut(_this.clearOverlay);
						$(_this.controls).addClass("loggedIn");
					} else {
						$("input[type=password]", _this.overlay).val("");
						$(".err", _this.overlay).text(resp.err);
					}
				});
			});
			$(_this.overlay).fadeIn();
		});
	};

	DesktopUI.prototype.showConfirmation = function(text, words, cb) {
		if (cb === undefined) {
			cb = words;
			words = undefined;
		}
		var _this = this;
		$(_this.overlay.firstChild).load("templates/desktop/confirm.htm", function() {
			$(_this.overlay).fadeIn();
			if (words !== undefined) {
				$(".ok", _this.overlay).val(words[0]);
				$(".cancel", _this.overlay).val(words[1]);
			}
			var resp = false;
			$("h1", _this.overlay).text(text);
			$(".ok", _this.overlay).click(function() {
				resp = true;
				respond();
			});
			$(".cancel", _this.overlay).click(function() {
				resp = false;
				respond();
			});

			function respond() {
				$(_this.overlay).fadeOut(_this.clearOverlay);
				cb(resp);
			}
		});
	};

	DesktopUI.prototype.clearOverlay = function() {
		$(this.firstChild).empty();
	};

	DesktopUI.prototype.slideInControls = function() {
		adjustWidth(this, 30);
	};

	DesktopUI.prototype.slideOutControls = function() {
		adjustWidth(this, 300);
	};

	function adjustWidth(ui, width) {
		$("#month", ui.container).animate({"padding-right": width+"px"});
		$("#name", ui.container).animate({"margin-right": width+"px"});
		$(ui.controls).animate({"width": width+"px"}, {"progress": function(){$(window).trigger("resize.month");}});
	}

	global.DesktopUI = DesktopUI;

	var monthNames = ["January","February","March","April","May","June",
			"July","August","September","October","November","December"],
		dayNames = ["Sunday","Monday","Tuesday","Wednesday","Thursday",
			"Friday","Saturday"];

	function initControls(ui, model, root) {
		model.getLoggedIn(function(loggedIn) {
			if (loggedIn) {
				$(root).addClass("loggedIn");
			}
		});

		$(root).load("templates/desktop/controls.htm", function() {
			$("#navigation .left", this).click(function() {
				prevMonth(ui);
			});
			$("#navigation .right", this).click(function() {
				nextMonth(ui);
			});
			$("#loginIndicator", this).click(function() {
				model.getLoggedIn(function(loggedIn) {
					if (loggedIn) {
						ui.showConfirmation("Really logout?", ["Logout", "Cancel"], function(doLogout) {
							if (doLogout) {
								model.logout(function() {
									$(root).removeClass("loggedIn");
								});
							}
						});
					} else {
						ui.displayLoginForm();
					}
				});
			});
		});
	}

	function initContainer(ui, root) {
		var changeMonth = function(e) {
			e = e.originalEvent;
			if (e.wheelDelta > 0) {
				prevMonth(ui);
			} else if (e.wheelDelta < 0) {
				nextMonth(ui);
			}
		};

		$(root).bind("mousewheel", function(e) {
			if (!ui.transitioning) {
				changeMonth(e);
			}
		});
	}

	function initOverlay(ui, root) {
		root.appendChild(document.createElement("div"));
		root.firstChild.className = "content";
		$(root).click(function(e) {
			if (e.target === root) {
				$(this).fadeOut();
			}
		});
	}

	function nextMonth(ui) {
		ui.currentDate.addMonths(1, true);
		ui.selectedDate = null;
		ui.updateDisplay();
	}

	function prevMonth(ui) {
		ui.currentDate.addMonths(-1, true);
		ui.selectedDate = null;
		ui.updateDisplay();
	}

	function displayControls(cal, root) {
	}

	function displayMonth(root, ui, events) {
		var fragment = document.createDocumentFragment(),
			calendar = $("<table />", {'id':'month'}),
			month = $("<div />", {'id':'name'}),
			year = $("<div />", {'id':'year'}),
			cont = $(root),
			row = $("<tr />"),
			date = ui.currentDate;

		dayNames.forEach(function(day) {
			row.append($("<th />").text(day.toLowerCase()));
		});
		calendar.append(row);

		$(fragment).append([calendar, month, year]);
		month.text(monthNames[date.getMonth()].toLowerCase());
		year.text(date.getFullYear());

		var firstVisibleDate = date.clone()
								   .setDate(1);
		firstVisibleDate.addDays(-firstVisibleDate.getDay());

		var lastVisibleDate = date.clone()
								  .setDate(1)
								  .addMonths(1)
								  .addDays(-1);

		var iterDate = firstVisibleDate.clone().clearTime();

		// We quit when we pass the last day of the month and it is a Sunday.
		while (iterDate <= lastVisibleDate || iterDate.getDay() !== 0) {

			// On the first day of the week, start a new row.
			if (iterDate.getDay() === 0) {
				row = $("<tr />");
				calendar.append(row);
			}

			var cell = $("<td />");
			var cellNum = $("<div class='date' />");
			cellNum.text(iterDate.toString("dd"));
			cell.data("date", iterDate.clone()).append(cellNum);
			row.append(cell);

			if (iterDate.getMonth() !== date.getMonth()) {
				cell.addClass("sideMonth");
			} else {
				if (iterDate.valueOf() === ui.model.getDate().valueOf()) {
					cell.addClass("today");
				}
				if (ui.selectedDate !== null && 
					iterDate.valueOf() === ui.selectedDate.valueOf()) {
					cell.append($("<div class='selectTriangle' />"));
				}
				$.each(events, function(ignored, e) {
					var happensOn = new XDate(e.startTime);
					if (happensOn.clone().clearTime().valueOf() === iterDate.valueOf()) {
						var eventDiv = $("<div class='event' />");
						var duration = Math.round(e.duration / 1000 / 1000 / 1000 / 60); // Convert to Minutes
						eventDiv.addClass(duration + "min");
						eventDiv.text(e.name);
						cell.append(eventDiv);
					}
				});
			}

			iterDate.addDays(1);
		}

		// Move the selected indicator around when people click on dates.
		$(calendar).delegate("td:not(.sideMonth)", "click", function(e) {
			var _this = this;
			var triangle = $(".selectTriangle", calendar);
			if (triangle.length === 0) {
				triangle = $("<div class='selectTriangle' />");
				$(this).append(triangle);
				triangle.fadeIn(250);
				ui.selectedDate = $(_this).data("date");
			} else {
				triangle.fadeOut(100, function() {
					$(_this).append(this);
					$(this).fadeIn(250);
					ui.selectedDate = $(_this).data("date");
				});
			}
		});

		// Remove the event handler when we're gone for efficiency.
		$(window).off("resize.month");
		cont.empty().append(fragment);

		var cal = $(calendar), mon = $(month);

		var defaultFontSize = mon.css("font-size").replace("px", "");

		var fixFont = function() {
			var currentSize = mon.css("font-size").replace("px","");
			if (mon.width() > (cal.width() - 60)) {
				var scaleFontRatio = (mon.width() / (cal.width() - 60));
				var scalePosRatio = (mon.width() / (cal.width()));
				var newSize = currentSize / scaleFontRatio;
				mon.css("font-size",  newSize + "px");
				mon.css("bottom", scalePosRatio * cont.height() * mon.css("bottom").replace("%", "") / 100 + "px");
			} else if (defaultFontSize > currentSize) {
				mon.css("font-size", defaultFontSize + "px");
				fixFont();
			}
		};

		// Fix the font's size while we fade in.
		var fixFontInterval = setInterval(fixFont, 10);
		setTimeout(function() {
			clearInterval(fixFontInterval);
		}, 200);

		// Fix it again when the window is resized.
		$(window).on("resize.month", fixFont);
	};
})(window);
(function(global) {
	"use strict";
	var visibleApi = {
		available: false
	};

	(function() {
		var hidden, state, visibilityChange; 
		if (typeof document.hidden !== "undefined") {
			hidden = "hidden";
			visibilityChange = "visibilitychange";
			state = "visibilityState";
		} else if (typeof document.mozHidden !== "undefined") {
			hidden = "mozHidden";
			visibilityChange = "mozvisibilitychange";
			state = "mozVisibilityState";
		} else if (typeof document.msHidden !== "undefined") {
			hidden = "msHidden";
			visibilityChange = "msvisibilitychange";
			state = "msVisibilityState";
		} else if (typeof document.webkitHidden !== "undefined") {
			hidden = "webkitHidden";
			visibilityChange = "webkitvisibilitychange";
			state = "webkitVisibilityState";
		}

		if (hidden !== undefined) {
			visibleApi = {
				available: true,
				hidden: hidden,
				visibilityChange: visibilityChange,
				state: state
			};
		}
	}());

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

		//displayControls(calendar, this.controls);

		this.events = [];
		calendar.getEventsForMonth(this.currentDate.clone(), function(data) {
			if (data.err === undefined) {
				_this.events = data.events;
				displayMonth(_this.container, _this, data.events);

				expireOldEvents();
				checkForNextDay();
			}
		});

		var today = XDate.today();
		var battery = navigator.battery || navigator.webkitBattery || navigator.mozBattery;

		var eventsTimeout = null;
		var midnightTimeout = null;

		// Expire events that happened before now.
		var expireOldEvents = function() {
			var now = new XDate();

			$("td:not(.sideMonth) .event.future, td:not(.sideMonth) .event.active", _this.container).each(function(i, div) {
				var e = $(div).data("event");
				var startTime = new XDate(e.startTime, true).setUTCMode(false, true);
				var endTime = startTime.clone().addMilliseconds(e.duration / 1000000);

				// If the event has already started
				if (startTime.diffMinutes(now) >= 0) {
					$(div).removeClass("future");

					// We're now past the starting point. If we're not past the
					// finishing point, then it's happening right NOW!
					if (endTime.diffMinutes(now) < 0) {
						$(div).addClass("active");
					} else {
						// Otherwise, the event is complete.
						$(div).removeClass("active");
					}
				}
			});

			var timeout = 1000 * 10; // 10 seconds
			if (battery && !battery.charging) {
				timeout *= 6; // 1 minute
			}

			// If we're invisible, don't schedule a new event.
			if (!visibleApi.available || document[visibleApi.state] !== "hidden") {
				eventsTimeout = setTimeout(expireOldEvents, timeout);
			} else {
				eventsTimeout = null;
			}
		};

		var checkForNextDay = function() {
			if (XDate.today().valueOf() !== today.valueOf()) {
				today = XDate.today();
				$(".today", _this.container).removeClass("today");
				$("td", _this.container).each(function(i, td) {
					var date = $(td).data("date");
					if (date.valueOf() === XDate.today().valueOf()) {
						$(td).addClass("today");
						return false;
					}
				});
			}

			var timeout = 500;
			if (battery && !battery.charging) {
				timeout *= 4;
			}

			// If we're invisible, don't schedule a new event.
			if (!visibleApi.available || document[visibleApi.state] !== "hidden") {
				midnightTimeout = setTimeout(checkForNextDay, timeout);
			} else {
				midnightTimeout = null;
			}
		};

		// If we can check for visibility, then we need to restart the checker when we flip back to the calendar.
		if (visibleApi.available) {
			document.addEventListener(visibleApi.visibilityChange, function() {
				if (document[visibleApi.state] === "visible") {
					console.log("Visible again. Catching up.");
					checkForNextDay();
					expireOldEvents();
				} else if (document[visibleApi.state] === "hidden") {
					clearTimeout(eventsTimeout);
					clearTimeout(midnightTimeout);
				}
			});
		}

		checkForNextDay();
		expireOldEvents();
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
				_this.events = data.events;
				if (data.err === undefined) {
					displayMonth(_this.container, _this, data.events);
					$(window).trigger("resize");
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
						_this.updateDisplay();
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
		$("*", this.firstChild).remove();
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
		$(ui.controls).animate({"width": width+"px"}, {"progress": function(){$(window).trigger("resize");}});
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
			$("#navigation .typcn-chevron-left", this).click(function() {
				prevMonth(ui);
			});
			$("#navigation .typcn-chevron-right", this).click(function() {
				nextMonth(ui);
			});
			$("#loginIndicator", this).click(function() {
				model.getLoggedIn(function(loggedIn) {
					if (loggedIn) {
						ui.showConfirmation("Really logout?", ["Logout", "Cancel"], function(doLogout) {
							if (doLogout) {
								model.logout(function() {
									$(root).removeClass("loggedIn");
									ui.updateDisplay();
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

		$(root).on("click", ".event", function(e) {
			console.log(e);
			return false;
		});

		$(root).on("click", "#month td:not(.sideMonth)", function(e) {
			var _this = this;
			var triangle = $(".selectTriangle", root);
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
	}

	function initOverlay(ui, root) {
		root.appendChild(document.createElement("div"));
		root.firstChild.className = "content";
		$(root).click(function(e) {
			if (e.target === root) {
				$(this).fadeOut(ui.clearOverlay);
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

		var sidebarWidth = $("#month", root).css("padding-right");
		calendar.css("padding-right", sidebarWidth);
		month.css("margin-right", sidebarWidth);

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
			}
			// We want the cell to expand when it's too big.
			cell.mouseenter(function() {
				var td = $(this);
				var lastEvent = td.children(".event").last();
				if (lastEvent.length === 0) {
					return;
				}
				var targetHeight = lastEvent.css("top").replace("px", "") - 0 + lastEvent.outerHeight();
				if (targetHeight <= td.outerHeight()) {
					return;
				}
				var startingStyle = {
					"height": td.height() + "px",
					"padding-bottom": td.css("padding-bottom")
				};
				var offset = td.offset();
				td.stop().css({
					"top": offset.top,
					"left": offset.left,
					"width": td.width() + "px"
				}).animate({
					"height": targetHeight - (td.outerHeight() - td.height()),
					"padding-bottom": "0.5em"
				}, function() {
					$(this).removeClass("expanding");
				}).addClass("expanded");

				// Only save the current state if we weren't already
				// collapsing, cause those values will be where we were
				// mid-collapse.
				if (!$(this).hasClass("collapsing")) {
					td.data("startingStyle", startingStyle);
				}
			}).mouseleave(function() {
				var td = $(this);
				td.addClass("collapsing").stop().animate(td.data("startingStyle"), function() {
					$(this).removeClass("expanded collapsing")
						   .css({
							   "top": "",
							   "left": "",
							   "width": "",
							   "height": "",
							   "padding-bottom": ""
					});
				});
			});

			// Mark today.
			if (iterDate.valueOf() === XDate.today().valueOf()) {
				cell.addClass("today");
			}
			// Add the selected indicator.
			if (ui.selectedDate !== null && 
				iterDate.valueOf() === ui.selectedDate.valueOf()) {
				cell.append($("<div class='selectTriangle' />"));
			}
			var eventsOnDay = [];
			$.each(events, function(ignored, e) {
				var happensOn = new XDate(e.startTime);
				if (happensOn.clone().clearTime().valueOf() === iterDate.valueOf()) {
					eventsOnDay.push(e);
				}
			});

			eventsOnDay.sort(function(a, b) {
				return new XDate(b.startTime).diffMinutes(new XDate(a.startTime));
			});

			$.each(eventsOnDay, function(i, e) {
				var eventDiv = $("<div class='event'><span class='name' /><span class='time' /></div>");
				var duration = Math.round(e.duration / 1000 / 1000 / 1000 / 60); // Convert to Minutes
				eventDiv.addClass(duration + "min")
						.css({"top": (2.5*(i+1)) + "em"})
						.data("event", e);
				$(".name", eventDiv).text(e.name);
				$(".time", eventDiv).text(new XDate(e.startTime, true).toString("h(:mm)TT"));
				cell.append(eventDiv);

				var endTime = new XDate(e.startTime, true).setUTCMode(false, true).addMilliseconds(e.duration / 1000000);

				if (endTime.valueOf() > new XDate().valueOf()) {
					eventDiv.addClass("future");
				}

				ui.model.getCalendarColor(e.cid, function(color) {
					eventDiv.css({"border-color":"#"+color});
				});
			});

			iterDate.addDays(1);
		}

		// Remove the event handler when we're gone for efficiency.
		$(window).off("resize.month");
		cont.empty().append(fragment);

		var cal = $(calendar), mon = $(month);

		var fixFont = function() {
			var currentSize = mon.css("font-size").replace("px","");

			mon.css("font-size", "");
			var defaultFontSize = mon.css("font-size").replace("px", "");

			mon.css("font-size", currentSize + "px");

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

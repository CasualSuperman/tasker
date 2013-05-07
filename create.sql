CREATE TABLE IF NOT EXISTS Users (
	uid INT NOT NULL AUTO_INCREMENT,
	PRIMARY KEY(uid),

	# Passwords stored using the bcrypt hash, which has a maximum length of 60
	# characters.
	password CHAR(60) BINARY,
	# RFC 5321 specifies that emails have a maximum length of 255
	email VARCHAR(255),
	# What this user wants others to identify them as
	displayName VARCHAR(100),
	activated BOOLEAN
);

# Every calendar must be visible by at least one person, or we have no reason
# to store it.  This allows us to make the uid part of the primary key.
CREATE TABLE IF NOT EXISTS Calendars (
	cid INT NOT NULL AUTO_INCREMENT,
	owner INT NOT NULL,
	# Use the calendar id and the user id as a paired primary key.
	PRIMARY KEY(cid),

	# What this user named this calendar.
	name VARCHAR(200),

	# The color (as a hex code) that the user picked for this calendar.
	color CHAR(6),

	FOREIGN KEY(owner) REFERENCES Users(uid)
);

# Unique event ids allow for sharing.
CREATE TABLE IF NOT EXISTS Events (
	eid INT NOT NULL AUTO_INCREMENT,
	PRIMARY KEY(eid),

	creator INT NOT NULL,
	FOREIGN KEY (creator) REFERENCES Users(uid),
	calendar INT NOT NULL,
	FOREIGN KEY (calendar) REFERENCES Calendars(cid),

	name VARCHAR(100),

	# If the event lasts all day.
	allDay BOOLEAN,
	# When the event starts.
	# 2013-01-22 11:00
	# January 22, 2013 11am
	start CHAR(17),
	# When it ends.
	end CHAR(17),

	# How we repeat
	# 0: No repeating
	# 1: Daily repeat
	# 2: Weekly repeat
	# 3: Monthly repeat
	# 4: Yearly repeat
	repeatType TINYINT,

	# If/how often the event repeats.
	repeatFrequency INTEGER,
	# When it stops repeating.
	repeatUntil DATE,

	# -5 to 5, or null.
	# if NULL and repeatType is monthly, look at the startDate to find the date
	# we repeat on.
	weekOfMonth TINYINT,

	# Used like a bitmask, _SMTWTFS
	days TINYINT,

	# Only consider full weeks during the month.
	fullWeek BOOLEAN
);

# Do we do this on a per-user basis, or globally?
#CREATE TABLE IF NOT EXISTS ModifiedEvents (
#	meid INT NOT NULL,
#	eid INT NOT NULL,
#	FOREIGN KEY (eid) REFERENCES Events(eid) ON DELETE CASCADE,
#	PRIMARY KEY(eid,meid),
#
#	originalStartDate DATE,
#
#	newStartDate DATE,
#	newEndDate DATE,
#
#	newStartTime TIME,
#	newEndTime TIME
#);

# Each event can have multiple owners, each owner can have multiple events.
# (Many to Many, with full participation from events).
CREATE TABLE IF NOT EXISTS EventShares (
	uid INT NOT NULL,
	eid INT NOT NULL,
	cid INT NOT NULL,
	PRIMARY KEY (uid, eid),
	newName VARCHAR(100),

	# permissions INTEGER

	FOREIGN KEY(uid) REFERENCES Users(uid),
	FOREIGN KEY(eid) REFERENCES Events(eid),
	FOREIGN KEY(cid) REFERENCES Calendars(cid)
);

# One event can be in multiple calendars, one calendar can have multiple
# events. One calendar cannot have the same event twice. (Many to Many, with
# full participation from events).
CREATE TABLE IF NOT EXISTS CalendarShares (
	uid INT NOT NULL,
	cid INT NOT NULL,
	PRIMARY KEY (uid, cid),

	newCalendarName VARCHAR(200),

	FOREIGN KEY(uid) REFERENCES Users(uid),
	FOREIGN KEY(cid) REFERENCES Calendars(cid)
);

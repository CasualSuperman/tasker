CREATE TABLE Users (
	uid INT NOT NULL AUTO_INCREMENT,
	PRIMARY KEY(uid),

	# Passwords stored using the bcrypt hash, which has a maximum length of 60
	# characters.
	password CHAR(60) BINARY,
	# RFC 5321 specifies that emails have a maximum length of 255
	email VARCHAR(255),
	activated BOOLEAN
);

# Every calendar must be visible by at least one person, or we have no reason
# to store it.  This allows us to make the uid part of the primary key.
CREATE TABLE Calendars (
	cid INT NOT NULL,
	uid INT NOT NULL,
	# Use the calendar id and the user id as a paired primary key.
	PRIMARY KEY(uid,cid),

	# If this user can share this calendar with others.
	canShare BOOLEAN,
	# If this user can write to this calendar.
	canWrite BOOLEAN,
	# What this user named this calendar.
	name VARCHAR(200),
	FOREIGN KEY(uid) REFERENCES Users
);

# Unique event ids allow for sharing.
CREATE TABLE Events (
	eid INT NOT NULL,
	PRIMARY KEY(eid),

	# If the event lasts all day.
	allDay BOOLEAN,
	# When the even starts.
	startTime TIME,
	# When it ends.
	endTime TIME,
	# The starting date.
	startDate DATE,
	# The ending date.
	endDate DATE,
	# If the event repeats.
	repeats BOOLEAN,
	# When it stops repeating.
	repeatUntil DATE,
	# In what fashion the even repeats
	# 0: By day of week
	# 1: By date of month
	# 
	repeatType INT,
	# A string that varies in format by repeatType ID
	repeatData VARCHAR(7)
);

# Each event can have multiple owners, each owner can have multiple events.
# (Many to Many, with full participation from events).
CREATE TABLE EventOwners (
	uid INT NOT NULL,
	eid INT NOT NULL,
	PRIMARY KEY (uid, eid),
	FOREIGN KEY(uid) REFERENCES Users,
	FOREIGN KEY(eid) REFERENCES Events
);

# One event can be in multiple calendars, one calendar can have multiple
# events. One calendar cannot have the same event twice. (Many to Many, with
# full participation from events).
CREATE TABLE CalendarEvents (
	eid INT NOT NULL,
	cid INT NOT NULL,
	PRIMARY KEY (eid, cid),
	FOREIGN KEY(eid) REFERENCES Events,
	FOREIGN KEY(cid) REFERENCES Calendars
);

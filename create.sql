PRAGMA synchronous = OFF;
PRAGMA journal_mode = MEMORY;
BEGIN TRANSACTION;
CREATE TABLE "CalendarShares" (
  "uid" INTEGER NOT NULL,
  "cid" INTEGER NOT NULL,
  "newCalendarName" varchar(200) DEFAULT NULL,
  PRIMARY KEY ("uid","cid")
  CONSTRAINT "CalendarShares_ibfk_1" FOREIGN KEY ("uid") REFERENCES "Users" ("uid"),
  CONSTRAINT "CalendarShares_ibfk_2" FOREIGN KEY ("cid") REFERENCES "Calendars" ("cid")
);
CREATE TABLE "Calendars" (
  "cid" INTEGER NOT NULL PRIMARY KEY,
  "owner" INTEGER NOT NULL,
  "name" varchar(200) DEFAULT NULL,
  "color" char(6) DEFAULT NULL,
  CONSTRAINT "Calendars_ibfk_1" FOREIGN KEY ("owner") REFERENCES "Users" ("uid")
);
CREATE TABLE "EventShares" (
  "uid" INTEGER NOT NULL,
  "eid" INTEGER NOT NULL,
  "cid" INTEGER NOT NULL,
  "newName" varchar(100) DEFAULT NULL,
  PRIMARY KEY ("uid","eid")
  CONSTRAINT "EventShares_ibfk_1" FOREIGN KEY ("uid") REFERENCES "Users" ("uid"),
  CONSTRAINT "EventShares_ibfk_2" FOREIGN KEY ("eid") REFERENCES "Events" ("eid"),
  CONSTRAINT "EventShares_ibfk_3" FOREIGN KEY ("cid") REFERENCES "Calendars" ("cid")
);
CREATE TABLE "Events" (
  "eid" INTEGER NOT NULL PRIMARY KEY,
  "creator" INTEGER NOT NULL,
  "calendar" INTEGER NOT NULL,
  "name" varchar(100) DEFAULT NULL,
  "allDay" tinyint(1) DEFAULT NULL,
  "repeatType" tinyint(4) DEFAULT NULL,
  "repeatFrequency" INTEGER DEFAULT NULL,
  "repeatUntil" date DEFAULT NULL,
  "weekOfMonth" tinyint(4) DEFAULT NULL,
  "days" tinyint(4) DEFAULT NULL,
  "fullWeek" tinyint(1) DEFAULT NULL,
  "start" char(17) DEFAULT NULL,
  "end" char(17) DEFAULT NULL,
  CONSTRAINT "Events_ibfk_1" FOREIGN KEY ("creator") REFERENCES "Users" ("uid"),
  CONSTRAINT "Events_ibfk_2" FOREIGN KEY ("calendar") REFERENCES "Calendars" ("cid")
);
CREATE TABLE "Users" (
  "uid" INTEGER NOT NULL PRIMARY KEY,
  "password" char(60) DEFAULT NULL,
  "email" varchar(255) DEFAULT NULL,
  "displayName" varchar(100) DEFAULT NULL,
  "activated" tinyint(1) DEFAULT NULL
);
CREATE TABLE Tasks (
	"tid" INTEGER NOT NULL PRIMARY KEY,
	"creator" INTEGER NOT NULL,
	"name" varchar(100) DEFAULT NULL,
	"timeRequired" INTEGER,
	"timeInvested" INTEGER NOT NULL,
	"dueDate" date DEFAULT NULL,
	CONSTRAINT "Tasks_ibfk_1" FOREIGN KEY ("creator") REFERENCES "Users" ("uid")
);
CREATE TABLE TaskInstances (
	"tiid" INTEGER NOT NULL PRIMARY KEY,
	"tid" INTEGER NOT NULL,
	"length" INTEGER NOT NULL,
	"when" date NOT NULL,
	"completed" tinyint(1) DEFAULT NULL,
	CONSTRAINT "TaskInstances_ibfk_1" FOREIGN KEY ("tid") REFERENCES "Tasks" ("tid")
);
CREATE TABLE EventInstances (
	eiid INTEGER NOT NULL,
	eid INTEGER NOT NULL,
	removed tinyint(1) NOT NULL DEFAULT "0",
	allDay tinyint(1) NOT NULL DEFAULT "0",
	start CHAR(17) DEFAULT NULL,
	end CHAR(17) DEFAULT NULL,
	CONSTRAINT "EventInstances_ibfk_1" FOREIGN KEY ("eid") REFERENCES "Events" ("eid"),
	PRIMARY KEY (eiid, eid)
);
CREATE INDEX "Calendars_owner" ON "Calendars" ("owner");
CREATE INDEX "Events_creator" ON "Events" ("creator");
CREATE INDEX "Events_calendar" ON "Events" ("calendar");
CREATE INDEX "EventShares_eid" ON "EventShares" ("eid");
CREATE INDEX "EventShares_cid" ON "EventShares" ("cid");
CREATE INDEX "CalendarShares_cid" ON "CalendarShares" ("cid");
END TRANSACTION;

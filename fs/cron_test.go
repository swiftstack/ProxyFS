package fs

import (
	"fmt"
	"testing"
	"time"

	"github.com/swiftstack/ProxyFS/conf"
)

func TestTimeParsing(t *testing.T) {
	tp, _ := time.Parse("09:30", "11:57")
	fmt.Println(tp)
	tn := time.Now()
	tz, off := tn.Zone()
	fmt.Println("tn  ==", tn)
	fmt.Println("tz  ==", tz)
	fmt.Println("off ==", off)
	loc := time.FixedZone("UTC-8", -8*60*60)
	td := time.Date(2009, time.November, 10, 23, 0, 0, 0, loc)
	fmt.Println("td in RFC822                format is:", td.Format(time.RFC822))
	fmt.Println("td in RFC3339               format is:", td.Format(time.RFC3339))
	fmt.Println("tn in RFC822                format is:", tn.Format(time.RFC822))
	fmt.Println("tn in RFC3339               format is:", tn.Format(time.RFC3339))
	fmt.Println("tn in \"2006-01-02T15:04:05\" format is:", tn.Format("2006-01-02T15:04:05"))
	ll, _ := time.LoadLocation("America/New_York")
	fmt.Println("time.LoadLocation(\"America/New_York\") ==", ll)
}

func TestLoadSnapShotPolicy(t *testing.T) {
	var (
		err                error
		snapShotPolicy     *snapShotPolicyStruct
		testConfMap        conf.ConfMap
		testConfMapStrings []string
	)

	// Case 0 - no SnapShotPolicy

	testConfMapStrings = []string{}

	testConfMap, err = conf.MakeConfMapFromStrings(testConfMapStrings)
	if nil != err {
		t.Fatalf("Case 0: conf.MakeConfMapFromStrings() failed: %v", err)
	}

	snapShotPolicy, err = loadSnapShotPolicy(testConfMap, "TestVolume")
	if nil != err {
		t.Fatalf("Case 0: loadSnapShotPolicy(testConfMap, \"TestVolume\") failed: %v", err)
	}

	if nil != snapShotPolicy {
		t.Fatalf("Case 0: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned non-nil snapShotPolicy")
	}

	// Case 1 - SnapShotPolicy with empty ScheduleList and no TimeZone

	testConfMapStrings = []string{
		"SnapShotPolicy:CommonSnapShotPolicy.ScheduleList=",
		"Volume:TestVolume.SnapShotPolicy=CommonSnapShotPolicy",
	}

	testConfMap, err = conf.MakeConfMapFromStrings(testConfMapStrings)
	if nil != err {
		t.Fatalf("Case 1: conf.MakeConfMapFromStrings() failed: %v", err)
	}

	snapShotPolicy, err = loadSnapShotPolicy(testConfMap, "TestVolume")
	if nil != err {
		t.Fatalf("Case 1: loadSnapShotPolicy(testConfMap, \"TestVolume\") failed: %v", err)
	}

	if "CommonSnapShotPolicy" != snapShotPolicy.name {
		t.Fatalf("Case 1: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy with unexpected .name")
	}
	if 0 != len(snapShotPolicy.schedule) {
		t.Fatalf("Case 1: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy with unexpected .schedule")
	}
	if "UTC" != snapShotPolicy.location.String() {
		t.Fatalf("Case 1: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy with unexpected .location")
	}

	// Case 2 - SnapShotPolicy with empty ScheduleList and empty TimeZone

	testConfMapStrings = []string{
		"SnapShotPolicy:CommonSnapShotPolicy.ScheduleList=",
		"SnapShotPolicy:CommonSnapShotPolicy.TimeZone=",
		"Volume:TestVolume.SnapShotPolicy=CommonSnapShotPolicy",
	}

	testConfMap, err = conf.MakeConfMapFromStrings(testConfMapStrings)
	if nil != err {
		t.Fatalf("Case 2: conf.MakeConfMapFromStrings() failed: %v", err)
	}

	snapShotPolicy, err = loadSnapShotPolicy(testConfMap, "TestVolume")
	if nil != err {
		t.Fatalf("Case 2: loadSnapShotPolicy(testConfMap, \"TestVolume\") failed: %v", err)
	}

	if "CommonSnapShotPolicy" != snapShotPolicy.name {
		t.Fatalf("Case 2: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy with unexpected .name")
	}
	if 0 != len(snapShotPolicy.schedule) {
		t.Fatalf("Case 2: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy with unexpected .schedule")
	}
	if "UTC" != snapShotPolicy.location.String() {
		t.Fatalf("Case 2: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy with unexpected .location")
	}

	// Case 3 - SnapShotPolicy with empty ScheduleList and TimeZone of "UTC"

	testConfMapStrings = []string{
		"SnapShotPolicy:CommonSnapShotPolicy.ScheduleList=",
		"SnapShotPolicy:CommonSnapShotPolicy.TimeZone=UTC",
		"Volume:TestVolume.SnapShotPolicy=CommonSnapShotPolicy",
	}

	testConfMap, err = conf.MakeConfMapFromStrings(testConfMapStrings)
	if nil != err {
		t.Fatalf("Case 3: conf.MakeConfMapFromStrings() failed: %v", err)
	}

	snapShotPolicy, err = loadSnapShotPolicy(testConfMap, "TestVolume")
	if nil != err {
		t.Fatalf("Case 3: loadSnapShotPolicy(testConfMap, \"TestVolume\") failed: %v", err)
	}

	if "CommonSnapShotPolicy" != snapShotPolicy.name {
		t.Fatalf("Case 3: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy with unexpected .name")
	}
	if 0 != len(snapShotPolicy.schedule) {
		t.Fatalf("Case 3: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy with unexpected .schedule")
	}
	if "UTC" != snapShotPolicy.location.String() {
		t.Fatalf("Case 3: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy with unexpected .location")
	}

	// Case 4 - SnapShotPolicy with empty ScheduleList and TimeZone of "Local"

	testConfMapStrings = []string{
		"SnapShotPolicy:CommonSnapShotPolicy.ScheduleList=",
		"SnapShotPolicy:CommonSnapShotPolicy.TimeZone=Local",
		"Volume:TestVolume.SnapShotPolicy=CommonSnapShotPolicy",
	}

	testConfMap, err = conf.MakeConfMapFromStrings(testConfMapStrings)
	if nil != err {
		t.Fatalf("Case 4: conf.MakeConfMapFromStrings() failed: %v", err)
	}

	snapShotPolicy, err = loadSnapShotPolicy(testConfMap, "TestVolume")
	if nil != err {
		t.Fatalf("Case 4: loadSnapShotPolicy(testConfMap, \"TestVolume\") failed: %v", err)
	}

	if "CommonSnapShotPolicy" != snapShotPolicy.name {
		t.Fatalf("Case 4: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy with unexpected .name")
	}
	if 0 != len(snapShotPolicy.schedule) {
		t.Fatalf("Case 4: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy with unexpected .schedule")
	}
	if "Local" != snapShotPolicy.location.String() {
		t.Fatalf("Case 4: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy with unexpected .location")
	}

	// Case 5 - SnapShotPolicy with exhaustive ScheduleList and a specific TimeZone

	testConfMapStrings = []string{
		"SnapShotSchedule:MinutelySnapShotSchedule.CronTab=* * * * *",
		"SnapShotSchedule:MinutelySnapShotSchedule.Keep=59",
		"SnapShotSchedule:HourlySnapShotSchedule.CronTab=0 * * * *",
		"SnapShotSchedule:HourlySnapShotSchedule.Keep=23",
		"SnapShotSchedule:DailySnapShotSchedule.CronTab=0 0 * * *",
		"SnapShotSchedule:DailySnapShotSchedule.Keep=6",
		"SnapShotSchedule:WeeklySnapShotSchedule.CronTab=0 0 * * 0",
		"SnapShotSchedule:WeeklySnapShotSchedule.Keep=8",
		"SnapShotSchedule:MonthlySnapShotSchedule.CronTab=0 0 1 * *",
		"SnapShotSchedule:MonthlySnapShotSchedule.Keep=11",
		"SnapShotSchedule:YearlySnapShotSchedule.CronTab=0 0 1 1 *",
		"SnapShotSchedule:YearlySnapShotSchedule.Keep=4",
		"SnapShotPolicy:CommonSnapShotPolicy.ScheduleList=MinutelySnapShotSchedule,HourlySnapShotSchedule,DailySnapShotSchedule,WeeklySnapShotSchedule,MonthlySnapShotSchedule,YearlySnapShotSchedule",
		"SnapShotPolicy:CommonSnapShotPolicy.TimeZone=America/Los_Angeles",
		"Volume:TestVolume.SnapShotPolicy=CommonSnapShotPolicy",
	}

	testConfMap, err = conf.MakeConfMapFromStrings(testConfMapStrings)
	if nil != err {
		t.Fatalf("Case 5: conf.MakeConfMapFromStrings() failed: %v", err)
	}

	snapShotPolicy, err = loadSnapShotPolicy(testConfMap, "TestVolume")
	if nil != err {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") failed: %v", err)
	}

	if "CommonSnapShotPolicy" != snapShotPolicy.name {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy with unexpected .name")
	}

	if 6 != len(snapShotPolicy.schedule) {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy with unexpected .schedule")
	}

	if "MinutelySnapShotSchedule" != snapShotPolicy.schedule[0].name {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[0] with unexpected .name")
	}
	if snapShotPolicy.schedule[0].minuteSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[0] with unexpected .minuteSpecified")
	}
	if snapShotPolicy.schedule[0].hourSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[0] with unexpected .hourSpecified")
	}
	if snapShotPolicy.schedule[0].dayOfMonthSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[0] with unexpected .dayOfMonthSpecified")
	}
	if snapShotPolicy.schedule[0].monthSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[0] with unexpected .monthSpecified")
	}
	if snapShotPolicy.schedule[0].dayOfWeekSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[0] with unexpected .dayOfWeekSpecified")
	}
	if 59 != snapShotPolicy.schedule[0].keep {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[0] with unexpected .keep")
	}

	if "HourlySnapShotSchedule" != snapShotPolicy.schedule[1].name {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[1] with unexpected .name")
	}
	if !snapShotPolicy.schedule[1].minuteSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[1] with unexpected .minuteSpecified")
	}
	if 0 != snapShotPolicy.schedule[1].minute {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[1] with unexpected .minute")
	}
	if snapShotPolicy.schedule[1].hourSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[1] with unexpected .hourSpecified")
	}
	if snapShotPolicy.schedule[1].dayOfMonthSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[1] with unexpected .dayOfMonthSpecified")
	}
	if snapShotPolicy.schedule[1].monthSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[1] with unexpected .monthSpecified")
	}
	if snapShotPolicy.schedule[1].dayOfWeekSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[1] with unexpected .dayOfWeekSpecified")
	}
	if 23 != snapShotPolicy.schedule[1].keep {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[1] with unexpected .keep")
	}

	if "DailySnapShotSchedule" != snapShotPolicy.schedule[2].name {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[2] with unexpected .name")
	}
	if !snapShotPolicy.schedule[2].minuteSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[2] with unexpected .minuteSpecified")
	}
	if 0 != snapShotPolicy.schedule[2].minute {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[2] with unexpected .minute")
	}
	if !snapShotPolicy.schedule[2].hourSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[2] with unexpected .hourSpecified")
	}
	if 0 != snapShotPolicy.schedule[2].hour {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[2] with unexpected .hour")
	}
	if snapShotPolicy.schedule[2].dayOfMonthSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[2] with unexpected .dayOfMonthSpecified")
	}
	if snapShotPolicy.schedule[2].monthSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[2] with unexpected .monthSpecified")
	}
	if snapShotPolicy.schedule[2].dayOfWeekSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[2] with unexpected .dayOfWeekSpecified")
	}
	if 6 != snapShotPolicy.schedule[2].keep {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[2] with unexpected .keep")
	}

	if "WeeklySnapShotSchedule" != snapShotPolicy.schedule[3].name {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[3] with unexpected .name")
	}
	if !snapShotPolicy.schedule[3].minuteSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[3] with unexpected .minuteSpecified")
	}
	if 0 != snapShotPolicy.schedule[3].minute {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[3] with unexpected .minute")
	}
	if !snapShotPolicy.schedule[3].hourSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[3] with unexpected .hourSpecified")
	}
	if 0 != snapShotPolicy.schedule[3].hour {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[3] with unexpected .hour")
	}
	if snapShotPolicy.schedule[3].dayOfMonthSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[3] with unexpected .dayOfMonthSpecified")
	}
	if snapShotPolicy.schedule[3].monthSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[3] with unexpected .monthSpecified")
	}
	if !snapShotPolicy.schedule[3].dayOfWeekSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[3] with unexpected .dayOfWeekSpecified")
	}
	if 0 != snapShotPolicy.schedule[3].dayOfWeek {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[3] with unexpected .dayOfWeek")
	}
	if 8 != snapShotPolicy.schedule[3].keep {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[3] with unexpected .keep")
	}

	if "MonthlySnapShotSchedule" != snapShotPolicy.schedule[4].name {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[4] with unexpected .name")
	}
	if !snapShotPolicy.schedule[4].minuteSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[4] with unexpected .minuteSpecified")
	}
	if 0 != snapShotPolicy.schedule[4].minute {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[4] with unexpected .minute")
	}
	if !snapShotPolicy.schedule[4].hourSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[4] with unexpected .hourSpecified")
	}
	if 0 != snapShotPolicy.schedule[4].hour {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[4] with unexpected .hour")
	}
	if !snapShotPolicy.schedule[4].dayOfMonthSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[4] with unexpected .dayOfMonthSpecified")
	}
	if 1 != snapShotPolicy.schedule[4].dayOfMonth {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[4] with unexpected .dayOfMonth")
	}
	if snapShotPolicy.schedule[4].monthSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[4] with unexpected .monthSpecified")
	}
	if snapShotPolicy.schedule[4].dayOfWeekSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[4] with unexpected .dayOfWeekSpecified")
	}
	if 11 != snapShotPolicy.schedule[4].keep {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[4] with unexpected .keep")
	}

	if "YearlySnapShotSchedule" != snapShotPolicy.schedule[5].name {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[5] with unexpected .name")
	}
	if !snapShotPolicy.schedule[5].minuteSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[5] with unexpected .minuteSpecified")
	}
	if 0 != snapShotPolicy.schedule[5].minute {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[5] with unexpected .minute")
	}
	if !snapShotPolicy.schedule[5].hourSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[5] with unexpected .hourSpecified")
	}
	if 0 != snapShotPolicy.schedule[5].hour {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[5] with unexpected .hour")
	}
	if !snapShotPolicy.schedule[5].dayOfMonthSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[5] with unexpected .dayOfMonthSpecified")
	}
	if 1 != snapShotPolicy.schedule[5].dayOfMonth {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[5] with unexpected .dayOfMonth")
	}
	if !snapShotPolicy.schedule[5].monthSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[5] with unexpected .monthSpecified")
	}
	if 1 != snapShotPolicy.schedule[5].month {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[5] with unexpected .month")
	}
	if snapShotPolicy.schedule[5].dayOfWeekSpecified {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[5] with unexpected .dayOfWeekSpecified")
	}
	if 4 != snapShotPolicy.schedule[5].keep {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy.schedule[5] with unexpected .keep")
	}

	if "America/Los_Angeles" != snapShotPolicy.location.String() {
		t.Fatalf("Case 5: loadSnapShotPolicy(testConfMap, \"TestVolume\") returned snapShotPolicy with unexpected .location")
	}
}

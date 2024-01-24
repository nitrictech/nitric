// Copyright Nitric Pty Ltd.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deploy

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/robfig/cron/v3"
)

// This file has mostly copied from
// https://github.com/aws/copilot-cli/blob/677dda6cddc7ff0c6ef60452ca40e2a528588c1e/internal/pkg/deploy/cloudformation/stack/scheduled_job.go
// The reasons for doing this are:
// - the functions that we need are all private.
// - it does not handle */x => 0/x

const (
	// Cron expressions in AWS Cloudwatch are of the form "M H DoM Mo DoW Y"
	// We use these predefined schedules when a customer specifies "@daily" or "@annually"
	// to fulfill the predefined schedules spec defined at
	// https://godoc.org/github.com/robfig/cron#hdr-Predefined_schedules
	// AWS requires that cron expressions use a ? wildcard for either DoM or DoW
	// so we represent that here.
	//            M H mD Mo wD Y
	cronHourly  = "0 * * * ? *" // at minute 0
	cronDaily   = "0 0 * * ? *" // at midnight
	cronWeekly  = "0 0 ? * 1 *" // at midnight on sunday
	cronMonthly = "0 0 1 * ? *" // at midnight on the first of the month
	cronYearly  = "0 0 1 1 ? *" // at midnight on January 1

	hourly   = "@hourly"
	daily    = "@daily"
	midnight = "@midnight"
	weekly   = "@weekly"
	monthly  = "@monthly"
	yearly   = "@yearly"
	annually = "@annually"

	every = "@every "

	fmtRateScheduleExpression = "rate(%d %s)" // rate({duration} {units})
	fmtCronScheduleExpression = "cron(%s)"
)

var awsScheduleRegexp = regexp.MustCompile(`(?:rate|cron)\(.*\)`) // Validates that an expression is of the form rate(xyz) or cron(abc)

// ConvertToAWS converts the Schedule string to the format required by Cloudwatch Events
// https://docs.aws.amazon.com/lambda/latest/dg/services-cloudwatchevents-expressions.html
// Cron expressions must have an sixth "year" field, and must contain at least one ? (either-or)
// in either day-of-month or day-of-week.
// Day-of-week expressions are zero-indexed in Golang but one-indexed in AWS.
// @every cron definition strings are converted to rates.
// All others become cron expressions.
// Exception is made for strings of the form "rate( )" or "cron( )". These are accepted as-is and
// validated server-side by CloudFormation.

// cron(Minutes Hours Day-of-month Month Day-of-week Year)

func ConvertToAWS(schedule string) (string, error) {
	if schedule == "" {
		return "", errors.New("schedule can not be empty")
	}

	// If the schedule uses default CloudWatch Events syntax, pass it through for server-side validation.
	if match := awsScheduleRegexp.FindStringSubmatch(schedule); match != nil {
		return schedule, nil
	}

	// Try parsing the string as a cron expression to validate it.
	if _, err := cron.ParseStandard(schedule); err != nil {
		return "", fmt.Errorf("schedule is not valid cron, rate, or preset: %w", err)
	}

	var (
		scheduleExpression string
		err                error
	)

	switch {
	case strings.HasPrefix(schedule, every):
		scheduleExpression, err = toRate(schedule[len(every):])
		if err != nil {
			return "", fmt.Errorf("parse fixed interval: %w", err)
		}
	case strings.HasPrefix(schedule, "@"):
		scheduleExpression, err = toFixedSchedule(schedule)
		if err != nil {
			return "", fmt.Errorf("parse preset schedule: %w", err)
		}
	default:
		scheduleExpression, err = toAWSCron(schedule)
		if err != nil {
			return "", fmt.Errorf("parse cron schedule: %w", err)
		}
	}

	return scheduleExpression, nil
}

// toRate converts a cron "@every" directive to a rate expression defined in minutes.
// example input: @every 1h30m
//
//	output: rate(90 minutes)
func toRate(duration string) (string, error) {
	d, err := time.ParseDuration(duration)
	if err != nil {
		return "", fmt.Errorf("parse duration: %w", err)
	}
	// Check that rates are not specified in units smaller than minutes
	if d != d.Truncate(time.Minute) {
		return "", fmt.Errorf("duration must be a whole number of minutes or hours")
	}

	if d < time.Minute*1 {
		return "", errors.New("duration must be greater than or equal to 1 minute")
	}

	minutes := int(d.Minutes())
	if minutes == 1 {
		return fmt.Sprintf(fmtRateScheduleExpression, minutes, "minute"), nil
	}

	return fmt.Sprintf(fmtRateScheduleExpression, minutes, "minutes"), nil
}

// toFixedSchedule converts cron predefined schedules into AWS-flavored cron expressions.
// (https://godoc.org/github.com/robfig/cron#hdr-Predefined_schedules)
// Example input: @daily
//
//	output: cron(0 0 * * ? *)
//	 input: @annually
//	output: cron(0 0 1 1 ? *)
func toFixedSchedule(schedule string) (string, error) {
	switch {
	case strings.HasPrefix(schedule, hourly):
		return fmt.Sprintf(fmtCronScheduleExpression, cronHourly), nil
	case strings.HasPrefix(schedule, midnight):
		fallthrough
	case strings.HasPrefix(schedule, daily):
		return fmt.Sprintf(fmtCronScheduleExpression, cronDaily), nil
	case strings.HasPrefix(schedule, weekly):
		return fmt.Sprintf(fmtCronScheduleExpression, cronWeekly), nil
	case strings.HasPrefix(schedule, monthly):
		return fmt.Sprintf(fmtCronScheduleExpression, cronMonthly), nil
	case strings.HasPrefix(schedule, annually):
		fallthrough
	case strings.HasPrefix(schedule, yearly):
		return fmt.Sprintf(fmtCronScheduleExpression, cronYearly), nil
	default:
		return "", fmt.Errorf("unrecognized preset schedule %s", schedule)
	}
}

func awsCronFieldSpecified(input string) bool {
	return !strings.ContainsAny(input, "*?")
}

// toAWSCron converts "standard" 5-element crons into the AWS preferred syntax
// cron(* * * * ? *)
// MIN HOU DOM MON DOW YEA
// EITHER DOM or DOW must be specified as ? (either-or operator)
// BOTH DOM and DOW cannot be specified
// DOW numbers run 1-7, not 0-6
// Example input: 0 9 * * 1-5 (at 9 am, Monday-Friday)
//
//	: cron(0 9 ? * 2-6 *) (adds required ? operator, increments DOW to 1-index, adds year)
func toAWSCron(schedule string) (string, error) {
	const (
		MIN = iota
		HOU
		DOM
		MON
		DOW
	)

	// Split the cron into its components. We can do this because it'll already have been validated.
	// Use https://golang.org/pkg/strings/#Fields since it handles consecutive whitespace.
	sched := strings.Fields(schedule)

	// Check whether the Day of Week and Day of Month fields have a ?
	// Possible conversion:
	// * * * * * ==> * * * * ?
	// 0 9 * * 1 ==> 0 9 ? * 1
	// 0 9 1 * * ==> 0 9 1 * ?
	switch {
	// If both are unspecified, convert DOW to a ? and DOM to *
	case !awsCronFieldSpecified(sched[DOM]) && !awsCronFieldSpecified(sched[DOW]):
		sched[DOW] = "?"
		sched[DOM] = "*"
	// If DOM is * or ? and DOW is specified, convert DOM to ?
	case !awsCronFieldSpecified(sched[DOM]) && awsCronFieldSpecified(sched[DOW]):
		sched[DOM] = "?"
	// If DOW is * or ? and DOM is specified, convert DOW to ?
	case !awsCronFieldSpecified(sched[DOW]) && awsCronFieldSpecified(sched[DOM]):
		sched[DOW] = "?"
	// Error if both DOM and DOW are specified
	default:
		return "", errors.New("cannot specify both DOW and DOM in cron expression")
	}

	// We also need to adjust the Day of Week value
	// crontab uses 0-6
	// AWS uses 1-7 (Sunday-Saturday)
	var newDOW []rune

	for _, c := range sched[DOW] {
		if unicode.IsDigit(c) {
			// Check for standard 0-6 day range and increment
			newDOW = append(newDOW, c+1)
		} else {
			newDOW = append(newDOW, c)
		}
	}
	// We don't need to use a string builder here because this will only ever have a max of
	// about 50 characters (SUN-MON,MON-TUE,TUE-WED,... is the longest possible string here)
	sched[DOW] = string(newDOW)

	// Add "every year" to 5-element crons to comply with AWS
	sched = append(sched, "*")

	return fmt.Sprintf(fmtCronScheduleExpression, strings.Join(sched, " ")), nil
}

// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import scheduler "Server/scheduler"
import time "time"

// IScheduler is an autogenerated mock type for the IScheduler type
type IScheduler struct {
	mock.Mock
}

// AddDailyAvailability provides a mock function with given fields: hour, minute, duration
func (_m *IScheduler) AddDailyAvailability(hour int, minute int, duration time.Duration) {
	_m.Called(hour, minute, duration)
}

// AddWeeklyAvailability provides a mock function with given fields: dayOfWeek, hour, minute, duration
func (_m *IScheduler) AddWeeklyAvailability(dayOfWeek time.Weekday, hour int, minute int, duration time.Duration) {
	_m.Called(dayOfWeek, hour, minute, duration)
}

// GetSchedule provides a mock function with given fields:
func (_m *IScheduler) GetSchedule() []scheduler.Game {
	ret := _m.Called()

	var r0 []scheduler.Game
	if rf, ok := ret.Get(0).(func() []scheduler.Game); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]scheduler.Game)
		}
	}

	return r0
}

// InitScheduler provides a mock function with given fields: tournamentType, roundsPerWeek, concurrentGameNum, gameDuration, start, end, teams
func (_m *IScheduler) InitScheduler(tournamentType int, roundsPerWeek int, concurrentGameNum int, gameDuration time.Duration, start time.Time, end time.Time, teams []int) {
	_m.Called(tournamentType, roundsPerWeek, concurrentGameNum, gameDuration, start, end, teams)
}

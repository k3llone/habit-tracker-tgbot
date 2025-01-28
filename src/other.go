package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func CheckHabitToday(userId int64, habitId int64, db *sql.DB) bool {
	today := time.Now()
	check_time := fmt.Sprintf("%v.%v.%v", today.Day(), today.Month(), today.Year())

	habit_complete := HabitComplete{}
	err := habit_complete.LoadDate(check_time, habitId, db)

	log.Println(err)

	if habit_complete.HabitId == habitId {
		return true
	}

	return false
}

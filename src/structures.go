package main

import (
	"database/sql"
	"log"
)

type User struct {
	Id          int64
	State       string
	CreateHabit int64
}

type Habit struct {
	Id      int64
	UserId  int64
	Name    string
	RemTime string
}

type HabitComplete struct {
	Id      int64
	HabitId int64
	Date    string
}

// USER

func (u *User) Insert(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO Users (id, states, createhabit) VALUES ($1, $2, $3)", u.Id, u.State, u.CreateHabit)

	if err != nil {
		return err
	}

	return nil
}

func (u *User) Update(db *sql.DB) error {
	_, err := db.Exec("UPDATE Users SET states = $1, createhabit = $2 WHERE id = $3", u.State, u.CreateHabit, u.Id)

	if err != nil {
		return err
	}

	return nil
}

func (u *User) Load(Id int64, db *sql.DB) error {
	res, err := db.Query("SELECT * FROM Users WHERE id=$1", Id)

	if err != nil {
		return err
	}

	var id int64
	state_str := ""

	defer res.Close()

	for res.Next() {
		res.Scan(&id, &state_str, &u.CreateHabit)
	}

	u.Id = id
	u.State = state_str

	return nil
}

// HABIT

func (h *Habit) Insert(db *sql.DB) error {
	res, err := db.Exec("INSERT INTO Habits (user, name, remtime) VALUES ($1, $2, $3)", h.UserId, h.Name, h.RemTime)

	if err != nil {
		log.Println(err)
	}

	h.Id, _ = res.LastInsertId()

	if err != nil {
		return err
	}

	return nil
}

func (h *Habit) Update(db *sql.DB) error {
	_, err := db.Exec("UPDATE Habits SET user = $1, name = $2, remtime = $3 WHERE id = $4", h.UserId, h.Name, h.RemTime, h.Id)

	if err != nil {
		return err
	}

	return nil
}

func (h *Habit) Load(Id int64, db *sql.DB) error {
	res, err := db.Query("SELECT * FROM Habits WHERE id=$1", Id)

	if err != nil {
		return err
	}

	defer res.Close()

	for res.Next() {
		res.Scan(&h.Id, &h.UserId, &h.Name, &h.RemTime)
	}

	return err
}

// HABIT_COMPLETE

func (hc *HabitComplete) Insert(db *sql.DB) error {
	res, err := db.Exec("INSERT INTO HabitComplete (habit, date) VALUES ($1, $2)", hc.HabitId, hc.Date)

	hc.Id, _ = res.LastInsertId()

	if err != nil {
		return err
	}

	return nil
}

func (hc *HabitComplete) Update(db *sql.DB) error {
	_, err := db.Exec("UPDATE HabitComplete SET habit = $1, date = $2 WHERE id = $3", hc.HabitId, hc.Date, hc.Id)

	if err != nil {
		return err
	}

	return nil
}

func (hc *HabitComplete) Load(Id int64, db *sql.DB) error {
	res, err := db.Query("SELECT * FROM HabitComplete WHERE id=$1", Id)

	if err != nil {
		return err
	}

	defer res.Close()

	for res.Next() {
		res.Scan(&hc.Id, &hc.HabitId, &hc.Date)
	}

	return err
}

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
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

type HabitMenu struct {
	Id     int64
	UserId int64
	Pages  int64
	Cpage  int64
	Habits []int64
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

func (h *Habit) Delete(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM Habits WHERE id=$1", h.Id)

	if err != nil {
		return err
	}

	return nil
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

func (hc *HabitComplete) LoadDate(Date string, HabitId int64, db *sql.DB) error {
	res, err := db.Query("SELECT * FROM HabitComplete WHERE date=$1 AND habit = $2", Date, HabitId)

	if res == nil {
		return errors.New("aaaa")
	}

	defer res.Close()

	for res.Next() {
		res.Scan(&hc.Id, &hc.HabitId, &hc.Date)
	}

	return err
}

// HABIT_MENU

func (hm *HabitMenu) Insert(db *sql.DB) error {
	habit_text := ""

	for i, v := range hm.Habits {
		if i == 0 {
			habit_text += fmt.Sprintf("%v", v)
		} else {
			habit_text += fmt.Sprintf(" %v", v)
		}
	}

	res, err := db.Exec("INSERT INTO HabitMenus (user, pages, cpage, habits) VALUES ($1, $2, $3, $4)", hm.UserId, hm.Pages, hm.Cpage, habit_text)

	hm.Id, _ = res.LastInsertId()

	if err != nil {
		return err
	}

	return nil
}

func (hm *HabitMenu) Update(db *sql.DB) error {
	_, err := db.Exec("UPDATE HabitMenus SET cpage = $1 WHERE id = $2", hm.Cpage, hm.Id)

	if err != nil {
		return err
	}

	return nil
}

func (hm *HabitMenu) Load(Id int64, db *sql.DB) error {
	res, err := db.Query("SELECT * FROM HabitMenus WHERE id=$1", Id)

	if err != nil {
		return err
	}

	defer res.Close()

	habits_text := ""

	for res.Next() {
		res.Scan(&hm.Id, &hm.UserId, &hm.Pages, &hm.Cpage, &habits_text)
	}

	habits := strings.Split(habits_text, " ")
	hm.Habits = make([]int64, 0)

	for _, v := range habits {
		num, _ := strconv.ParseInt(v, 10, 64)
		hm.Habits = append(hm.Habits, num)
	}

	return err
}

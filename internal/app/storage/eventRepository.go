package storage

import (
	"calendar/internal/app/model"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

// EventRepository for managing Event data
type EventRepository struct {
	store *Storage
}

func (r *EventRepository) Create(e *model.Event) (*model.Event, error) {
	preparedNotes := strings.Join(e.Notes, ",")
	if err := r.store.db.QueryRow(
		"INSERT INTO event (title, description, time, timezone, duration, notes) VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
		e.Title,
		e.Description,
		e.Time,
		e.Timezone,
		e.Duration,
		preparedNotes,
	).Scan(&e.ID); err != nil {
		return nil, err
	}

	return e, nil
}

func (r *EventRepository) Load(id int) (*model.Event, error) {
	savedNotes := ""
	e := &model.Event{}

	if err := r.store.db.QueryRow(
		"SELECT id, title, description, time, timezone, duration, notes FROM event WHERE id = $1",
		id,
	).Scan(&e.ID, &e.Title, &e.Description, &e.Time, &e.Timezone, &e.Duration, &savedNotes); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	e.Notes = strings.Split(savedNotes, ",")

	return e, nil
}

func (r *EventRepository) Update(event *model.Event) (*model.Event, error) {
	preparedNotes := strings.Join(event.Notes, ", ")
	e := &model.Event{}

	var currentId int = 0
	var savedTime, savedTimezone string
	if err := r.store.db.QueryRow(
		"SELECT id, time, timezone FROM event WHERE id = $1",
		event.ID,
	).Scan(&currentId, &savedTime, &savedTimezone); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	//layout := "2021-09-24 13:45:39.724447928 +0300 EEST m=+3.338521868"
	//layout := "2006-01-02T15:04:05.000Z"
	str := savedTime//"2014-11-12T11:45:26.371Z"
	t, erro := time.Parse(time.RFC3339, str) //time.Parse(layout, str)

	if erro != nil {
		fmt.Println(erro)
	}
	fmt.Println("PARSED TIME = ", t)
	if currentId > 0 {
		if _, err := r.store.db.Query(
			"UPDATE event SET title = $1, description = $2, time = $3, timezone = $4, duration = $5, notes = $6 WHERE id = $7",
			event.Title, event.Description, event.Time, event.Timezone, event.Duration, preparedNotes, event.ID,
		); err != nil {

			return nil, err
		}

		return e, nil
	}

	err := errors.New("not existing ID")
	return nil, err
}


func (r *EventRepository) All() ([]model.Event, error) {
	rows, err := r.store.db.Query(
		"SELECT id, title, description, time, timezone, duration, notes FROM event",
	)

	defer rows.Close()
	if err == sql.ErrNoRows {
		return nil, nil
	}

	var events []model.Event

	event := model.Event{}
	savedNote := ""
	for rows.Next() {
		if err := rows.Scan(&event.ID, &event.Title, &event.Description, &event.Time, &event.Timezone, &event.Duration, &savedNote); err != nil {
			log.Fatal(err)
		}

		event.Notes = strings.Split(savedNote, ", ")
		events = append(events, event)
		//log.Printf("id %d has role %s\n", event.ID, event)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return events, nil
}

// Delete event from db by id
func (r *EventRepository) Delete(id int) (bool, error) {
	result, err := r.store.db.Exec("DELETE FROM event WHERE id = $1", id)

	if err != nil {
		fmt.Println("deleting caused error")
		return false, err
	}
	fmt.Println("Deleted", result)

	return true, nil
}



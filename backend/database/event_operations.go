package database

import (
	"database/sql"
	"errors"
	"event_management/backend/models"
	"log"
)

func GetEventsByOrganizerID(organizerID int) ([]models.EventWithRegistrationCount, error) {
	events := []models.EventWithRegistrationCount{}

	rows, err := DB.Query(`
		SELECT 
			e.event_id, e.title, e.description, e.date, e.location, e.max_capacity, e.organiser_id, 
			CASE WHEN e.isalive = 0 THEN 'cancelled' ELSE 'active' END AS status,
			COUNT(r.registration_id) as registered_count
		FROM event e
		LEFT JOIN registration r ON e.event_id = r.event_id AND r.isalive = 1
		WHERE e.organiser_id = ? AND e.isalive = 1
		GROUP BY e.event_id
		ORDER BY e.date DESC
	`, organizerID)

	if err != nil {
		return events, err
	}
	defer rows.Close()

	for rows.Next() {
		var event models.EventWithRegistrationCount
		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Description,
			&event.Date,
			&event.Location,
			&event.Capacity,
			&event.OrganizerID,
			&event.Status,
			&event.RegisteredCount,
		)
		if err != nil {
			return events, err
		}
		events = append(events, event)
	}

	return events, nil
}

func GetRegistrationsByEventID(eventID int) ([]models.RegistrationWithUserDetails, error) {
	registrations := []models.RegistrationWithUserDetails{}

	rows, err := DB.Query(`
		SELECT 
			r.registration_id, r.event_id, r.attendee_id, r.registration_date, r.status,
			a.name, a.email
		FROM registration r
		JOIN attendee a ON r.attendee_id = a.attendee_id
		WHERE r.event_id = ? AND r.isalive = 1
		ORDER BY r.registration_date DESC
	`, eventID)

	if err != nil {
		return registrations, err
	}
	defer rows.Close()

	for rows.Next() {
		var reg models.RegistrationWithUserDetails
		err := rows.Scan(
			&reg.ID,
			&reg.EventID,
			&reg.UserID,
			&reg.RegistrationDate,
			&reg.Status,
			&reg.UserName,
			&reg.Email,
		)
		if err != nil {
			return registrations, err
		}
		registrations = append(registrations, reg)
	}

	return registrations, nil
}

func IsEventOrganizer(eventID, userID int) (bool, error) {
	var organizerID int

	err := DB.QueryRow("SELECT organiser_id FROM event WHERE event_id = ?", eventID).Scan(&organizerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("event not found")
		}
		return false, err
	}

	return organizerID == userID, nil
}

func CreateEvent(event models.Event) (models.Event, error) {
	if event.Status == "" {
		event.Status = "active"
	}

	isActive := true
	if event.Status == "cancelled" {
		isActive = false
	}

	log.Printf("Creating event with capacity: %d", event.Capacity)

	result, err := DB.Exec(
		"INSERT INTO event (title, description, date, location, max_capacity, organiser_id, isalive) VALUES (?, ?, ?, ?, ?, ?, ?)",
		event.Name, event.Description, event.Date, event.Location, event.Capacity, event.OrganizerID, isActive,
	)

	if err != nil {
		return event, err
	}

	eventID, err := result.LastInsertId()
	if err != nil {
		return event, err
	}

	var savedCapacity int
	err = DB.QueryRow("SELECT max_capacity FROM event WHERE event_id = ?", eventID).Scan(&savedCapacity)
	if err != nil {
		log.Printf("Warning: Could not verify saved capacity: %v", err)
	} else {
		log.Printf("Saved event capacity: %d", savedCapacity)
		event.Capacity = savedCapacity
	}

	event.ID = int(eventID)
	return event, nil
}

func UpdateEvent(event models.Event) (*models.Event, error) {
	if event.OrganizerID == 0 {
		return nil, errors.New("organizer ID is required for update authorization")
	}

	query := `UPDATE event
			  SET title = ?, description = ?, date = ?, location = ?, max_capacity = ?
			  WHERE event_id = ? AND organiser_id = ?`

	result, err := DB.Exec(query,
		event.Name,
		event.Description,
		event.Date,
		event.Location,
		event.Capacity,
		event.ID,
		event.OrganizerID,
	)

	if err != nil {
		log.Printf("Error executing update event statement for event %d: %v", event.ID, err)
		return nil, errors.New("failed to update event")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected for update event %d: %v", event.ID, err)
		return nil, errors.New("failed to confirm event update")
	}

	if rowsAffected == 0 {
		var exists bool
		checkQuery := "SELECT EXISTS(SELECT 1 FROM event WHERE event_id = ?)"
		errCheck := DB.QueryRow(checkQuery, event.ID).Scan(&exists)
		if errCheck != nil && errCheck != sql.ErrNoRows {
			log.Printf("Error checking event existence during update for event %d: %v", event.ID, errCheck)
			return nil, errors.New("failed to update event check")
		}
		if !exists {
			return nil, errors.New("event not found")
		} else {
			return nil, errors.New("unauthorized or event not updatable")
		}
	}

	return &event, nil
}

func CancelEvent(eventID int) error {
	_, err := DB.Exec("UPDATE event SET isalive = 0 WHERE event_id = ?", eventID)
	return err
}

func GetAllEvents() ([]models.EventWithRegistrationCount, error) {
	events := []models.EventWithRegistrationCount{}

	rows, err := DB.Query(`
		SELECT 
			e.event_id, e.title, e.description, e.date, e.location, e.max_capacity, e.organiser_id, 
			CASE WHEN e.isalive = 0 THEN 'cancelled' ELSE 'active' END AS status,
			COUNT(r.registration_id) as registered_count
		FROM event e
		LEFT JOIN registration r ON e.event_id = r.event_id AND r.isalive = 1
		WHERE e.isalive = 1
		GROUP BY e.event_id
		ORDER BY e.date ASC
	`)

	if err != nil {
		return events, err
	}
	defer rows.Close()

	for rows.Next() {
		var event models.EventWithRegistrationCount
		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Description,
			&event.Date,
			&event.Location,
			&event.Capacity,
			&event.OrganizerID,
			&event.Status,
			&event.RegisteredCount,
		)
		if err != nil {
			return events, err
		}
		events = append(events, event)
	}

	return events, nil
}

func GetEventByID(eventID int) (models.Event, error) {
	var event models.Event

	row := DB.QueryRow(`
		SELECT 
			event_id, title, description, date, location, max_capacity, organiser_id, 
			CASE WHEN isalive = 0 THEN 'cancelled' ELSE 'active' END AS status
		FROM event
		WHERE event_id = ? AND isalive = 1
	`, eventID)

	err := row.Scan(
		&event.ID,
		&event.Name,
		&event.Description,
		&event.Date,
		&event.Location,
		&event.Capacity,
		&event.OrganizerID,
		&event.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return event, errors.New("event not found")
		}
		return event, err
	}

	return event, nil
}

func IsUserRegisteredForEvent(userID, eventID int) (bool, error) {
	var count int
	err := DB.QueryRow(`
		SELECT COUNT(*) 
		FROM registration 
		WHERE attendee_id = ? AND event_id = ? AND isalive = 1
	`, userID, eventID).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func CreateRegistration(reg models.Registration) (int, error) {
	var capacity, registered int
	err := DB.QueryRow(`
		SELECT e.max_capacity, COUNT(r.registration_id) 
		FROM event e
		LEFT JOIN registration r ON e.event_id = r.event_id AND r.isalive = 1
		WHERE e.event_id = ? AND e.isalive = 1
		GROUP BY e.event_id
	`, reg.EventID).Scan(&capacity, &registered)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("event not found")
		}
		return 0, err
	}

	if registered >= capacity {
		return 0, errors.New("event is full")
	}

	result, err := DB.Exec(`
		INSERT INTO registration (event_id, attendee_id, registration_date, status, isalive)
		VALUES (?, ?, ?, ?, 1)
	`, reg.EventID, reg.UserID, reg.RegistrationDate, reg.Status)

	if err != nil {
		return 0, err
	}

	regID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(regID), nil
}

func IsRegistrationOwner(regID, userID int) (bool, error) {
	var count int
	err := DB.QueryRow(`
		SELECT COUNT(*) 
		FROM registration 
		WHERE registration_id = ? AND attendee_id = ? AND isalive = 1
	`, regID, userID).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func CancelRegistration(regID int) error {
	_, err := DB.Exec("UPDATE registration SET isalive = 0 WHERE registration_id = ?", regID)
	return err
}

package database

import (
	"database/sql"
	"errors"
	"event_management/backend/models"
	"log"
)

// GetEventsByOrganizerID returns all events created by a given organiser (user with role='organiser'),
// along with a count of active registrations.
func GetEventsByOrganizerID(organizerID int) ([]models.EventWithRegistrationCount, error) {
	events := []models.EventWithRegistrationCount{}

	rows, err := DB.Query(`
		SELECT 
			e.event_id, e.title, e.description, e.date, e.location, e.max_capacity, e.organiser_id, 
			CASE WHEN e.isalive = 0 THEN 'cancelled' ELSE 'active' END AS status,
			COUNT(r.registration_id) as registered_count
		FROM event e
		LEFT JOIN registration r 
		  ON e.event_id = r.event_id 
		 AND r.isalive = 1
		WHERE e.organiser_id = ? 
		  AND e.isalive = 1
		GROUP BY e.event_id
		ORDER BY e.date DESC
	`, organizerID)
	if err != nil {
		return events, err
	}
	defer rows.Close()

	for rows.Next() {
		var ev models.EventWithRegistrationCount
		if err := rows.Scan(
			&ev.ID,
			&ev.Name,
			&ev.Description,
			&ev.Date,
			&ev.Location,
			&ev.Capacity,
			&ev.OrganizerID,
			&ev.Status,
			&ev.RegisteredCount,
		); err != nil {
			return events, err
		}
		events = append(events, ev)
	}

	return events, rows.Err()
}

// GetRegistrationsByEventID returns all registrations for a given event, 
// along with the attendee’s user name and email pulled from the unified user table.
func GetRegistrationsByEventID(eventID int) ([]models.RegistrationWithUserDetails, error) {
	registrations := []models.RegistrationWithUserDetails{}

	rows, err := DB.Query(`
		SELECT 
			r.registration_id,
			r.event_id,
			r.attendee_id,
			r.registration_date,
			r.status,
			u.name,
			u.email
		FROM registration r
		JOIN user u 
		  ON r.attendee_id = u.user_id 
		 AND u.role = 'attendee'
		WHERE r.event_id = ? 
		  AND r.isalive = 1
		ORDER BY r.registration_date DESC
	`, eventID)
	if err != nil {
		return registrations, err
	}
	defer rows.Close()

	for rows.Next() {
		var reg models.RegistrationWithUserDetails
		if err := rows.Scan(
			&reg.ID,
			&reg.EventID,
			&reg.UserID,
			&reg.RegistrationDate,
			&reg.Status,
			&reg.UserName,
			&reg.Email,
		); err != nil {
			return registrations, err
		}
		registrations = append(registrations, reg)
	}

	return registrations, rows.Err()
}

// IsEventOrganizer checks whether the given user (by user_id) is the organiser of that event.
func IsEventOrganizer(eventID, userID int) (bool, error) {
	var organiserID int
	err := DB.QueryRow(`
		SELECT organiser_id 
		FROM event 
		WHERE event_id = ?
	`, eventID).Scan(&organiserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("event not found")
		}
		return false, err
	}
	return organiserID == userID, nil
}

// CreateEvent inserts a new event (organized by an existing user with role='organiser').
func CreateEvent(event models.Event) (models.Event, error) {
	if event.Status == "" {
		event.Status = "active"
	}
	isActive := event.Status != "cancelled"

	log.Printf("Creating event with capacity: %d", event.Capacity)
	res, err := DB.Exec(`
		INSERT INTO event 
			(title, description, date, location, max_capacity, organiser_id, isalive)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, event.Name, event.Description, event.Date, event.Location, event.Capacity, event.OrganizerID, isActive)
	if err != nil {
		return event, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return event, err
	}
	event.ID = int(lastID)
	return event, nil
}

// UpdateEvent modifies an existing event, ensuring only the owner (organiser) can update it.
func UpdateEvent(event models.Event) (*models.Event, error) {
	if event.OrganizerID == 0 {
		return nil, errors.New("organizer ID is required for update authorization")
	}

	res, err := DB.Exec(`
		UPDATE event
		SET title = ?, description = ?, date = ?, location = ?, max_capacity = ?
		WHERE event_id = ? 
		  AND organiser_id = ?
	`, event.Name, event.Description, event.Date, event.Location, event.Capacity, event.ID, event.OrganizerID)
	if err != nil {
		log.Printf("Error updating event %d: %v", event.ID, err)
		return nil, errors.New("failed to update event")
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return nil, errors.New("failed to confirm event update")
	}
	if ra == 0 {
		return nil, errors.New("no permission or event not found")
	}

	return &event, nil
}

// CancelEvent marks an event as inactive (soft delete).
func CancelEvent(eventID int) error {
	_, err := DB.Exec(`
		UPDATE event 
		SET isalive = 0 
		WHERE event_id = ?
	`, eventID)
	return err
}

// GetAllEvents fetches every active event with its registration count.
func GetAllEvents() ([]models.EventWithRegistrationCount, error) {
	events := []models.EventWithRegistrationCount{}

	rows, err := DB.Query(`
		SELECT 
			e.event_id, e.title, e.description, e.date, e.location, e.max_capacity, e.organiser_id,
			CASE WHEN e.isalive = 0 THEN 'cancelled' ELSE 'active' END AS status,
			COUNT(r.registration_id) as registered_count
		FROM event e
		LEFT JOIN registration r 
		  ON e.event_id = r.event_id 
		 AND r.isalive = 1
		WHERE e.isalive = 1
		GROUP BY e.event_id
		ORDER BY e.date ASC
	`)
	if err != nil {
		return events, err
	}
	defer rows.Close()

	for rows.Next() {
		var ev models.EventWithRegistrationCount
		if err := rows.Scan(
			&ev.ID,
			&ev.Name,
			&ev.Description,
			&ev.Date,
			&ev.Location,
			&ev.Capacity,
			&ev.OrganizerID,
			&ev.Status,
			&ev.RegisteredCount,
		); err != nil {
			return events, err
		}
		events = append(events, ev)
	}

	return events, rows.Err()
}

// GetEventByID retrieves a single event if it’s still active.
func GetEventByID(eventID int) (models.Event, error) {
	var ev models.Event
	err := DB.QueryRow(`
		SELECT 
			event_id, title, description, date, location, max_capacity, organiser_id,
			CASE WHEN isalive = 0 THEN 'cancelled' ELSE 'active' END AS status
		FROM event
		WHERE event_id = ? 
		  AND isalive = 1
	`, eventID).Scan(
		&ev.ID,
		&ev.Name,
		&ev.Description,
		&ev.Date,
		&ev.Location,
		&ev.Capacity,
		&ev.OrganizerID,
		&ev.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return ev, errors.New("event not found")
		}
		return ev, err
	}
	return ev, nil
}

// IsUserRegisteredForEvent checks if a given user (attendee) is signed up for an event.
func IsUserRegisteredForEvent(userID, eventID int) (bool, error) {
	var count int
	err := DB.QueryRow(`
		SELECT COUNT(*) 
		FROM registration 
		WHERE attendee_id = ? 
		  AND event_id = ? 
		  AND isalive = 1
	`, userID, eventID).Scan(&count)
	return count > 0, err
}

// CreateRegistration attempts to register a user (attendee) for an event if there’s capacity.
func CreateRegistration(reg models.Registration) (int, error) {
	var capacity, registered int
	err := DB.QueryRow(`
		SELECT e.max_capacity, COUNT(r.registration_id)
		FROM event e
		LEFT JOIN registration r 
		  ON e.event_id = r.event_id 
		 AND r.isalive = 1
		WHERE e.event_id = ? 
		  AND e.isalive = 1
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

	res, err := DB.Exec(`
		INSERT INTO registration 
			(event_id, attendee_id, registration_date, status, isalive)
		VALUES (?, ?, ?, ?, 1)
	`, reg.EventID, reg.UserID, reg.RegistrationDate, reg.Status)
	if err != nil {
		return 0, err
	}
	lastID, err := res.LastInsertId()
	return int(lastID), err
}

// IsRegistrationOwner ensures a user can only operate on their own registration.
func IsRegistrationOwner(regID, userID int) (bool, error) {
	var cnt int
	err := DB.QueryRow(`
		SELECT COUNT(*) 
		FROM registration 
		WHERE registration_id = ? 
		  AND attendee_id = ? 
		  AND isalive = 1
	`, regID, userID).Scan(&cnt)
	return cnt > 0, err
}

// CancelRegistration soft-deletes a registration.
func CancelRegistration(regID int) error {
	_, err := DB.Exec(`
		UPDATE registration 
		SET isalive = 0 
		WHERE registration_id = ?
	`, regID)
	return err
}

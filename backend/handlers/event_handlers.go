package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"event_management/backend/database"
	"event_management/backend/models"
	"event_management/backend/utils"

	"github.com/gorilla/mux"
)

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

type eventRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Location    string `json:"location"`
	Capacity    int    `json:"capacity"`
}

func GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := database.GetAllEvents()
	if err != nil {
		http.Error(w, "Error retrieving events", http.StatusInternalServerError)
		log.Printf("Error retrieving events: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func RegisterForEventHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	eventIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing event ID", http.StatusBadRequest)
		return
	}

	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}


	if _, err := database.GetEventByID(eventID); err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	isRegistered, err := database.IsUserRegisteredForEvent(userID, eventID)
	if err != nil {
		http.Error(w, "Error checking registration status", http.StatusInternalServerError)
		return
	}

	if isRegistered {
		http.Error(w, "User is already registered for this event", http.StatusConflict)
		return
	}

	registration := models.Registration{
		UserID:           userID,
		EventID:          eventID,
		RegistrationDate: time.Now().Format("2006-01-02 15:04:05"),
		Status:           "confirmed",
	}

	registrationID, err := database.CreateRegistration(registration)
	if err != nil {
		if err.Error() == "event is full" {
			http.Error(w, "Event is at full capacity", http.StatusBadRequest)
		} else {
			http.Error(w, "Error creating registration", http.StatusInternalServerError)
			log.Printf("Error creating registration: %v", err)
		}
		return
	}

	registration.ID = registrationID

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(registration)
}

func CancelRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	regIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing registration ID", http.StatusBadRequest)
		return
	}

	regID, err := strconv.Atoi(regIDStr)
	if err != nil {
		http.Error(w, "Invalid registration ID", http.StatusBadRequest)
		return
	}

	isOwner, err := database.IsRegistrationOwner(regID, userID)
	if err != nil {
		http.Error(w, "Error verifying registration ownership", http.StatusInternalServerError)
		return
	}

	if !isOwner {
		http.Error(w, "Unauthorized to cancel this registration", http.StatusForbidden)
		return
	}

	if err := database.CancelRegistration(regID); err != nil {
		http.Error(w, "Error cancelling registration", http.StatusInternalServerError)
		log.Printf("Error cancelling registration: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Registration cancelled successfully",
	})
}

func GetOrganizerEventsHandler(w http.ResponseWriter, r *http.Request) {
	organizerID, ok := r.Context().Value(utils.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	events, err := database.GetEventsByOrganizerID(organizerID)
	if err != nil {
		http.Error(w, "Error retrieving events", http.StatusInternalServerError)
		log.Printf("Error retrieving events: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func GetEventRegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	organizerID, ok := r.Context().Value(utils.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	eventIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing event ID", http.StatusBadRequest)
		return
	}

	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	if isOwner, err := database.IsEventOrganizer(eventID, organizerID); err != nil {
		http.Error(w, "Error verifying event ownership", http.StatusInternalServerError)
		return
	} else if !isOwner {
		http.Error(w, "Unauthorized to access this event's registrations", http.StatusForbidden)
		return
	}

	registrations, err := database.GetRegistrationsByEventID(eventID)
	if err != nil {
		http.Error(w, "Error retrieving registrations", http.StatusInternalServerError)
		log.Printf("Error retrieving registrations: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(registrations); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	organizerID, ok := r.Context().Value(utils.UserIDKey).(int)
	if !ok {
		writeJSONError(w, "Organizer ID not found in context", http.StatusInternalServerError)
		return
	}

	var req eventRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Name == "" || req.Date == "" || req.Location == "" {
		writeJSONError(w, "Event name, date, and location are required", http.StatusBadRequest)
		return
	}

	eventDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		writeJSONError(w, "Invalid date format. Please use YYYY-MM-DD.", http.StatusBadRequest)
		return
	}

	event := models.Event{
		OrganizerID: organizerID,
		Name:        req.Name,
		Description: req.Description,
		Date:        eventDate.Format("2006-01-02"),
		Location:    req.Location,
		Capacity:    req.Capacity,
	}

	createdEvent, err := database.CreateEvent(event)
	if err != nil {
		log.Printf("Error creating event: %v", err)
		writeJSONError(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdEvent)
}

func UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		writeJSONError(w, "Event ID is required", http.StatusBadRequest)
		return
	}
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, "Invalid Event ID format", http.StatusBadRequest)
		return
	}

	organizerID, ok := r.Context().Value(utils.UserIDKey).(int)
	if !ok {
		writeJSONError(w, "Organizer ID not found in context", http.StatusInternalServerError)
		return
	}

	var req eventRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSONError(w, "Invalid request body format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Name == "" || req.Date == "" || req.Location == "" {
		writeJSONError(w, "Event name, date, and location are required", http.StatusBadRequest)
		return
	}

	eventDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		writeJSONError(w, "Invalid date format. Please use YYYY-MM-DD.", http.StatusBadRequest)
		return
	}

	event := models.Event{
		ID:          eventID,
		OrganizerID: organizerID,
		Name:        req.Name,
		Description: req.Description,
		Date:        eventDate.Format("2006-01-02"),
		Location:    req.Location,
		Capacity:    req.Capacity,
	}

	updatedEvent, err := database.UpdateEvent(event)
	if err != nil {
		if err.Error() == "event not found" {
			writeJSONError(w, "Event not found", http.StatusNotFound)
		} else if err.Error() == "unauthorized" {
			writeJSONError(w, "You do not have permission to update this event", http.StatusForbidden)
		} else {
			log.Printf("Error updating event %d: %v", eventID, err)
			writeJSONError(w, "Failed to update event", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedEvent)
}

func CancelEventHandler(w http.ResponseWriter, r *http.Request) {
	organizerID, ok := r.Context().Value(utils.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	eventIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing event ID", http.StatusBadRequest)
		return
	}

	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	if isOwner, err := database.IsEventOrganizer(eventID, organizerID); err != nil {
		http.Error(w, "Error verifying event ownership", http.StatusInternalServerError)
		return
	} else if !isOwner {
		http.Error(w, "Unauthorized to cancel this event", http.StatusForbidden)
		return
	}

	if err := database.CancelEvent(eventID); err != nil {
		http.Error(w, "Error cancelling event", http.StatusInternalServerError)
		log.Printf("Error cancelling event: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Event cancelled successfully",
	})
}

package database

import (
	"event_management/backend/models"
	"fmt"
)

type UserData struct {
	Name  string
	Email string
}

// Fetch users by role
func getUsersByRole(role string) ([]UserData, error) {
	query := "SELECT name, email FROM user WHERE isalive = 1 AND role = ?"

	rows, err := DB.Query(query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []UserData
	for rows.Next() {
		var user UserData
		if err := rows.Scan(&user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// These are wrappers calling getUsersByRole with appropriate roles
func GetAllAdmins() ([]UserData, error) {
	return getUsersByRole("admin")
}

func GetAllOrganisers() ([]UserData, error) {
	return getUsersByRole("organiser")
}

func GetAllAttendees() ([]UserData, error) {
	return getUsersByRole("attendee")
}

type UserWithRole struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func GetAllUserRoles() ([]UserWithRole, error) {
	var allUsers []UserWithRole

	roles := []string{"admin", "organiser", "attendee"}
	for _, role := range roles {
		users, err := getUsersByRole(role)
		if err != nil {
			return nil, err
		}
		for _, user := range users {
			allUsers = append(allUsers, UserWithRole{
				Name:  user.Name,
				Email: user.Email,
				Role:  role,
			})
		}
	}
	return allUsers, nil
}

func DeactivateUser(email string, role string) error {
	query := "UPDATE user SET isalive = 0 WHERE email = ? AND role = ?"

	result, err := DB.Exec(query, email, role)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with email %s and role %s", email, role)
	}

	return nil
}

func GetUserByID(userID int) (models.User, error) {
	var user models.User

	err := DB.QueryRow(`
		SELECT user_id, name, email, phone, password, role
		FROM user
		WHERE user_id = ? AND isalive = 1
	`, userID).Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Password, &user.Role)

	return user, err
}

func UpdateUser(user models.User) (models.User, error) {
	_, err := DB.Exec(`
		UPDATE user
		SET name = ?, phone = ?
		WHERE user_id = ? AND isalive = 1
	`, user.Name, user.Phone, user.ID)

	if err != nil {
		return user, err
	}

	return GetUserByID(user.ID)
}

func GetRegistrationsByUserID(userID int) ([]models.Registration, error) {
	var registrations []models.Registration

	rows, err := DB.Query(`
		SELECT registration_id, event_id, attendee_id, registration_date, status
		FROM registration
		WHERE attendee_id = ? AND isalive = 1
		ORDER BY registration_date DESC
	`, userID)
	if err != nil {
		return registrations, err
	}
	defer rows.Close()

	for rows.Next() {
		var reg models.Registration
		err := rows.Scan(
			&reg.ID,
			&reg.EventID,
			&reg.UserID,
			&reg.RegistrationDate,
			&reg.Status,
		)
		if err != nil {
			return registrations, err
		}
		registrations = append(registrations, reg)
	}

	return registrations, nil
}

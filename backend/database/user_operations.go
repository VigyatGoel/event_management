package database

import (
	"fmt"
)

type UserData struct {
	Name  string
	Email string
}

func GetAllAdmins() ([]UserData, error) {
	query := "SELECT name, email FROM admin where isalive=1"

	rows, err := DB.Query(query)
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

func GetAllOrganisers() ([]UserData, error) {
	query := "SELECT name, email FROM organiser where isalive=1"

	rows, err := DB.Query(query)
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

func GetAllAttendees() ([]UserData, error) {
	query := "SELECT name, email FROM attendee where isalive=1"

	rows, err := DB.Query(query)
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

type UserWithRole struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func GetAllUserRoles() ([]UserWithRole, error) {
	var allUsers []UserWithRole

	admins, err := GetAllAdmins()
	if err != nil {
		return nil, err
	}
	for _, admin := range admins {
		allUsers = append(allUsers, UserWithRole{
			Name:  admin.Name,
			Email: admin.Email,
			Role:  "admin",
		})
	}

	organisers, err := GetAllOrganisers()
	if err != nil {
		return nil, err
	}
	for _, organiser := range organisers {
		allUsers = append(allUsers, UserWithRole{
			Name:  organiser.Name,
			Email: organiser.Email,
			Role:  "organiser",
		})
	}

	attendees, err := GetAllAttendees()
	if err != nil {
		return nil, err
	}
	for _, attendee := range attendees {
		allUsers = append(allUsers, UserWithRole{
			Name:  attendee.Name,
			Email: attendee.Email,
			Role:  "attendee",
		})
	}

	return allUsers, nil
}

func DeactivateUser(email string, role string) error {
	var updateQuery string

	switch role {
	case "admin":
		updateQuery = "UPDATE admin SET isalive = 0 WHERE email = ?"
	case "organiser":
		updateQuery = "UPDATE organiser SET isalive = 0 WHERE email = ?"
	case "attendee":
		updateQuery = "UPDATE attendee SET isalive = 0 WHERE email = ?"
	default:
		return fmt.Errorf("invalid role: %s", role)
	}

	result, err := DB.Exec(updateQuery, email)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with email %s in role %s", email, role)
	}

	return nil
}

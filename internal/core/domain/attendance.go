package domain

import "time"

type Task struct {
	ID                string    `json:"id"`
	OrgID             string    `json:"org_id"`
	Title             string    `json:"title" validate:"required"`
	AssignedUserID    *string   `json:"assigned_user_id,omitempty" validate:"omitempty,uuid"`
	GeofencingEnabled bool      `json:"geofencing_enabled"`
	LocationName      string    `json:"location_name,omitempty"`
	Latitude          float64   `json:"latitude,omitempty"`
	Longitude         float64   `json:"longitude,omitempty"`
	RadiusMeters      int       `json:"radius_meters,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

type Attendance struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	OrgID        string     `json:"org_id"`
	TaskID       *string    `json:"task_id,omitempty"`
	CheckInTime  time.Time  `json:"check_in_time"`
	CheckOutTime *time.Time `json:"check_out_time,omitempty"`
	Status       string     `json:"status"` // PRESENT, LATE, ABSENT
	Type         string     `json:"type"`   // GENERAL, TASK
	ShiftApplied string     `json:"shift_applied,omitempty"`
	LocationLat  float64    `json:"location_lat"`
	LocationLong float64    `json:"location_long"`
	Note         string     `json:"note,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type CheckInRequest struct {
	OrganizationID string  `json:"organization_id" validate:"required,uuid"`
	TaskID         *string `json:"task_id" validate:"omitempty,uuid"`
	Latitude       float64 `json:"latitude" validate:"required"`
	Longitude      float64 `json:"longitude" validate:"required"`
	Note           string  `json:"note"`
}

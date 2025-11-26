package domain

import "time"

type Organization struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name" validate:"required"`
	Email               string    `json:"email" validate:"required,email"`
	DefaultLocationLat  float64   `json:"default_location_lat" validate:"required"`
	DefaultLocationLong float64   `json:"default_location_long" validate:"required"`
	CreatedAt           time.Time `json:"created_at"`
}

type OrganizationMember struct {
	ID        string    `json:"id"`
	OrgID     string    `json:"org_id"`
	UserID    string    `json:"user_id"`
	Role      string    `json:"role"` // OWNER, MANAGER, EMPLOYEE
	GroupID   *string   `json:"group_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type OrganizationMemberDetail struct {
	OrganizationMember
	User User `json:"user"`
}

type Shift struct {
	ID                 string    `json:"id"`
	OrgID              string    `json:"org_id"`
	Name               string    `json:"name" validate:"required"`
	StartTime          string    `json:"start_time" validate:"required"` // HH:MM
	EndTime            string    `json:"end_time" validate:"required"`   // HH:MM
	Timezone           string    `json:"timezone" validate:"required"`
	AllowedLateMinutes int       `json:"allowed_late_minutes"`
	WorkingDays        []string  `json:"working_days" validate:"required,min=1"`
	CreatedAt          time.Time `json:"created_at"`
}

type Group struct {
	ID        string    `json:"id"`
	OrgID     string    `json:"org_id"`
	Name      string    `json:"name" validate:"required"`
	ShiftID   *string   `json:"shift_id,omitempty" validate:"omitempty,uuid"`
	ManagerID *string   `json:"manager_id,omitempty" validate:"omitempty,uuid"`
	CreatedAt time.Time `json:"created_at"`
}

type InviteEmployeeRequest struct {
	Email   string  `json:"email" validate:"required,email"`
	Role    string  `json:"role" validate:"required,oneof=MANAGER EMPLOYEE"`
	GroupID *string `json:"group_id" validate:"omitempty,uuid"`
}

type AssignUserToGroupRequest struct {
	GroupID string `json:"group_id" validate:"required,uuid"`
}

type UpdateEmployeeRequest struct {
	Role    string  `json:"role" validate:"required,oneof=MANAGER EMPLOYEE"`
	GroupID *string `json:"group_id" validate:"omitempty,uuid"`
}

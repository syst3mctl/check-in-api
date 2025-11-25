package domain

type GroupPerformanceReport struct {
	GroupName         string `json:"group_name"`
	AssignedShift     string `json:"assigned_shift"`
	TotalLateCheckins int    `json:"total_late_checkins"`
	AttendanceRate    string `json:"attendance_rate"`
}

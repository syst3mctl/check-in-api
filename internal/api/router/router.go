package router

import (
	_ "github.com/syst3mctl/check-in-api/docs" // Swagger docs
	"github.com/syst3mctl/check-in-api/internal/api/handler"
	"github.com/syst3mctl/check-in-api/internal/api/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func New(authHandler *handler.AuthHandler, orgHandler *handler.OrgHandler, attendanceHandler *handler.AttendanceHandler, reportHandler *handler.ReportHandler, authMiddleware *middleware.AuthMiddleware) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Handle)

		// Organizations
		r.Post("/organizations", orgHandler.CreateOrganization)
		r.Get("/organizations", orgHandler.ListOrganizations)
		r.Get("/organizations/{org_id}", orgHandler.GetOrganization)
		r.Put("/organizations/{org_id}", orgHandler.UpdateOrganization)
		r.Delete("/organizations/{org_id}", orgHandler.DeleteOrganization)
		r.Post("/organizations/{org_id}/invitations", orgHandler.InviteEmployee)
		r.Post("/organizations/{org_id}/shifts", orgHandler.CreateShift)
		r.Post("/organizations/{org_id}/groups", orgHandler.CreateGroup)
		r.Put("/organizations/{org_id}/members/{user_id}", orgHandler.AssignUserToGroup)
		r.Get("/organizations/{org_id}/employees", orgHandler.GetEmployees)
		r.Put("/organizations/{org_id}/employees/{user_id}", orgHandler.UpdateEmployee)
		r.Delete("/organizations/{org_id}/employees/{user_id}", orgHandler.RemoveEmployee)

		// Tasks
		r.Post("/organizations/{org_id}/tasks", attendanceHandler.CreateTask)

		// Attendance
		r.Post("/attendance/check-in", attendanceHandler.CheckIn)
		r.Post("/attendance/check-out", attendanceHandler.CheckOut)

		// Reports
		r.Get("/organizations/{org_id}/reports/groups/{group_id}", reportHandler.GetGroupPerformance)
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	return r
}

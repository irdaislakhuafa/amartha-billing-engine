package rest

func (r *rest) RegisterRoutes() {
	// server health and testing purpose
	r.svr.Get("/ping", r.Ping)

	api := r.svr.Group("api")
	v1 := api.Group("/v1")
	{
		v1.Post("/auth/register", r.Register)
		v1.Post("/auth/login", r.Login)

		v1.Get("/loans", r.ListLoan)
		v1.Post("/loans", r.CreateLoan)
		v1.Put("/loans", r.UpdateLoan)
		v1.Get("/loans/:id", r.GetLoan)
		v1.Delete("/loans/:id/:is_deleted", r.DeleteLoan)

		v1.Get("/users", r.authJWT, r.ListUser)
		v1.Put("/users", r.authJWT, r.UpdateUser)
		v1.Get("/users/:id", r.authJWT, r.GetUser)
		v1.Delete("/users/:id/:is_deleted", r.authJWT, r.DeleteUser)

		v1.Post("/loan/transactions", r.authJWT, r.CreateLoanTransaction)
		v1.Get("/loan/transactions", r.authJWT, r.ListLoanTransaction)
		v1.Get("/loan/transaction/calculate/:user_id", r.authJWT, r.CalculateOutstandingLoanTransaction)
		v1.Post("/loan/transaction/pay", r.authJWT, r.PayLoanTransaction)
		v1.Post("/loan/transaction/schedule-delinquent", r.authJWT, r.ScheduleDelinquent)

		v1.Post("/settings", r.authJWT, r.CreateSetting)
		v1.Get("/settings", r.authJWT, r.ListSetting)
		v1.Put("/settings", r.authJWT, r.UpdateSetting)
		v1.Get("/settings/:id", r.authJWT, r.GetSetting)

	}
}

package rest

func (r *rest) RegisterRoutes() {
	// server health and testing purpose
	r.svr.Get("/ping", r.Ping)

	api := r.svr.Group("api")
	v1 := api.Group("/v1")
	{
		v1.Get("/loans", r.ListLoan)
		v1.Post("/loans", r.CreateLoan)
		v1.Put("/loans", r.UpdateLoan)
		v1.Get("/loans/:id", r.GetLoan)
		v1.Delete("/loans/:id/:is_deleted", r.DeleteLoan)
	}
}

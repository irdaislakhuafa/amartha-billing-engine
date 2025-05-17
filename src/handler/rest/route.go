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
		v1.Post("/loans", r.authJWT, r.CreateLoan)
		v1.Put("/loans", r.authJWT, r.UpdateLoan)
		v1.Get("/loans/:id", r.GetLoan)
		v1.Delete("/loans/:id/:is_deleted", r.authJWT, r.DeleteLoan)

		v1.Get("/users", r.authJWT, r.ListUser)
		v1.Put("/users", r.authJWT, r.UpdateUser)
		v1.Get("/users/:id", r.authJWT, r.GetUser)
		v1.Delete("/users/:id/:is_deleted", r.authJWT, r.DeleteUser)

	}
}

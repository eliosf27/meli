package main

// Middleware build the routes of the server
func (s *Server) Routes() {
	s.server.GET("/", s.dependencies.StatusController.Status)

	items := s.server.Group("/items")
	items.GET("/:item_id", s.dependencies.ItemController.Get)
}

// Hooks defines an interface for a lifecycle of container.
// Example of using the Zeus container with Hooks and ErrorSet:
//
//	c := zeus.New()
//
//	c.Provide(func(h Hooks) *sql.DB {
//	   db := &sql.DB{} // pseudo-implementation
//
//	   h.OnStart(func() error {
//	      return errors.New("Failed to authenticate")
//	   })
//
//	   h.OnStop(func() error {
//	      return errors.New("Failed to close the database")
//	   })
//
//	   return db
//	})
package hooks

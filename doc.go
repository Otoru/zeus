// Package zeus provides a lightweight dependency injection container for Go.
// It allows users to register factories and resolve dependencies at runtime.
//
// Basic usage:
//
//	c := zeus.New()
//	c.Provide(func() int { return 42 })
//	c.Run(func(i int) {
//	    fmt.Println(i) // Outputs: 42
//	})
package zeus // import "github.com/otoru/zeus"

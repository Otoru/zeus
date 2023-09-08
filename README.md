# 🌩 Zeus - Simple Dependency Injection Container

![GitHub](https://img.shields.io/github/license/otoru/zeus)
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/otoru/zeus)
[![codecov](https://codecov.io/gh/Otoru/zeus/graph/badge.svg?token=Yfkyp5NZsY)](https://codecov.io/gh/Otoru/zeus)

Zeus is a sleek and efficient dependency injection container for Go. Easily register "factories" (functions that create instances of types) and let zeus resolve those dependencies at runtime.

## 🌟 Features

Why using zeus?

### 🚀 Simple to Use

With a minimalist API, integrating zeus into any Go project is a breeze.

### 🔍 Dependency Resolution

Register your dependencies and let zeus handle the rest.

### ⚠️ Cyclic Dependency Detection

Zeus detects and reports cycles in your dependencies to prevent runtime errors.

### 🪝 Hooks

Zeus supports lifecycle hooks, allowing you to execute functions at the start and end of your application. This is especially useful for setups and teardowns, like establishing a database connection or gracefully shutting down services.

## 🚀 Getting Started

### Installation

```bash
go get -u github.com/otoru/zeus
```

### Register Dependencies

```go
package main

import "github.com/otoru/zeus"

c := zeus.New()

c.Provide(func() int {
  return 42
})

c.Provide(func(i int) string {
  return fmt.Sprintf("Number: %d", i) 
})
```

### Resolve & Run Functions

```go
err := c.Run(func(s string) error {
    fmt.Println(s) // Outputs: Number: 42
    return nil
})
```

### Using Hooks

Zeus allows you to register hooks that run at the start and end of your application. This is useful for setting up and tearing down resources.

```go
c := zeus.New()

// Servoce is a dummy service that depends on Hooks.
type Service struct{}

c.Provide(func(h zeus.Hooks) *Service {
    h.OnStart(func() error {
        fmt.Println("Starting up...")
        return nil
    })

    h.OnStop(func() error {
        fmt.Println("Shutting down...")
        return nil
    })
    return &Service{}
})

c.Run(func(s *Service) {
    fmt.Println("Main function running with the service!")
})

// Outputs:
// Starting up...
// Main function running with the service!
// Shutting down...

```

### Error Handling

Zeus uses `ErrorSet` to aggregate multiple errors. This is especially useful when multiple errors occur during the lifecycle of your application, such as during dependency resolution or hook execution.

An ErrorSet can be returned from the Run method. Here's how you can handle it:

```go
err := c.Run(func() { /* ... */ })
if es, ok := err.(*zeus.ErrorSet); ok {
    for _, e := range es.Errors() {
        fmt.Println(e)
    }
}
```

## 🤝 Contributing

Contributions are warmly welcomed! Please open a PR or an issue if you find any problems or have enhancement suggestions.

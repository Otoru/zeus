# 🌩 Zeus - Simple Dependency Injection Container

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

## 🤝 Contributing

Contributions are warmly welcomed! Please open a PR or an issue if you find any problems or have enhancement suggestions.

# protoc-gen-gofilter


<div style="text-align: center; padding: 20px;">
    <img src=./.doc/imgs/gofilter.jpg width="50%" />
</div>

[![Go Reference](https://pkg.go.dev/badge/github.com/hsuanshao/protoc-gen-gofilter.svg)](https://pkg.go.dev/github.com/hsuanshao/protoc-gen-gofilter)
[![Go Report Card](https://goreportcard.com/badge/github.com/hsuanshao/protoc-gen-gofilter)](https://goreportcard.com/badge/github.com/hsuanshao/protoc-gen-gofilter)

High-performance, Zero-Reflection Field-Level Permission Control for Go Protobuf.

`protoc-gen-gofilter` is a Protoc plugin designed to solve "Field-Level" granular permission control issues in gRPC/Protobuf services. It generates efficient filtering code via Code Generation, avoiding the performance overhead and type unsafety of traditional Go Reflection.

## ✨ Features

*   **🚀 Extreme Performance**: Implemented based on BitSet (bitmask), filtering operations require only O(1) bitwise operations, which is tens of times faster than reflection.
*   **🛡️ Zero Reflection**: All filtering logic is generated at Compile-time, with no runtime reflection overhead.
*   **🔒 Declarative Definition**: Define permission rules directly in `.proto` files, serving as a Single Source of Truth (SSOT).
*   **🧩 Decoupled Architecture**: Separates "Permission Decision" from "Permission Enforcement", perfectly supporting mixed RBAC and ABAC models.
*   **🔧 Easy Integration**: Generated code provides standard interfaces, easily integrated into gRPC Interceptors or Middleware.

## 📦 Installation

### 1. Install Compiler Plugin (CLI Tool)

```bash
go install github.com/hsuanshao/protoc-gen-gofilter/cmd/protoc-gen-gofilter@latest
```

### 2. Download Runtime Library

Import the dependency in your project:

```bash
go get github.com/hsuanshao/protoc-gen-gofilter
```ds

## 🚀 Quick Start

### Step 1: Define Proto

In your `.proto` file, import `filter.proto` and use the `(filter.apply)` option to tag fields that require permission control.

> 💡 **Tip**: For convenience, it is recommended to copy `pb/filter.proto` from this project to your `third_party` directory, or use vendoring.

**myapp.proto**

```proto
syntax = "proto3";
package myapp;

option go_package = "example.com/my-project/pb/myapp";

// 1. Import definition
import "github.com/hsuanshao/protoc-gen-gofilter/protos/filter/filter.proto"; 

message UserProfile {
  int64 id = 1;
  string name = 2;

  // 2. Tag permission: Only users with "user.email.read" permission can see this field
  string email = 3 [(filter.apply) = "user.email.read"];

  // Supports repeated and optional fields
  repeated string secrets = 4 [(filter.apply) = "admin.secrets.read"];
}
```

### Step 2: Generate Go Code

When running `protoc`, add the `--go-filter_out` parameter.

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-filter_out=. --go-filter_opt=paths=source_relative \
       --proto_path=. \
       myapp.proto
```

After execution, you will see `myapp_filter.pb.go` generated, containing the `FilterFields` method.

### Step 3: Use in Go

In your business logic or gRPC Interceptor, calculate user permissions and execute filtering.

```go
package main

import (
    "fmt"
    "github.com/hsuanshao/protoc-gen-gofilter/entity/filter"
    "example.com/my-project/pb/myapp"
)

func main() {
    // 1. Simulate data
    data := &myapp.UserProfile{
        Id:      1,
        Name:    "Alice",
        Email:   "alice@example.com",
        Secrets: []string{"top-secret"},
    }

    // 2. [Decision Layer] Calculate permission BitMask based on user role
    // In real scenarios, this is usually combined with a Policy Engine (e.g., OPA, Casbin, or DB rules)
    userMask := filter.NewBitSet()

    // Assume the user has permission to read Email but not Secrets
    // Registry.GetID returns the unique integer ID for the permission string
    if emailPermID, ok := filter.Registry.GetID("user.email.read"); ok {
        userMask.Set(emailPermID)
    }

    // 3. [Enforcement Layer] Execute filtering
    // FilterFields is an auto-generated method, extremely fast
    data.FilterFields(userMask)

    // 4. Verify results
    fmt.Printf("Email: %s\n", data.Email)     // Output: alice@example.com
    fmt.Printf("Secrets: %v\n", data.Secrets) // Output: [] (filtered to nil)
}
```

## 🏗️ Architecture Design

This tool uses a **Permission Flattening** strategy:

1.  **Compile Time**: The Plugin scans `.proto` files, extracts all permission strings (e.g., `"user.email.read"`), and generates an `init()` function to automatically register them with the global `Registry` at startup, obtaining unique `int` IDs.
2.  **Runtime**:
    *   **BitSet**: Uses a dynamically sized `[]uint64` to represent sets of permissions.
    *   **Check**: Generated code directly uses `mask.Has(int_id)` for bitwise checks. If the check fails, the field is directly set to its Zero Value (`nil`, `0`, `""`).

This design moves the overhead of string comparison to startup time, leaving only extremely efficient integer operations during request processing.

## 🤝 Contributing

Contributions via Issues or Pull Requests are welcome!

1.  Fork this project
2.  Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3.  Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4.  Push to the Branch (`git push origin feature/AmazingFeature`)
5.  Open a Pull Request

## 📄 License

Distributed under the MIT License. See LICENSE for more information.

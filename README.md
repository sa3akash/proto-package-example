# @sa3akash/proto

> Protobuf + ConnectRPC definitions for **Users** and **Posts** services — with generated TypeScript and Go code, field-level validation via protovalidate, and automatic NPM publishing via GitHub Actions.

[![NPM Version](https://img.shields.io/npm/v/@sa3akash/proto.svg)](https://www.npmjs.com/package/@sa3akash/proto)
[![CI](https://github.com/sa3akash/proto-package-example/actions/workflows/ci.yml/badge.svg)](https://github.com/sa3akash/proto-package-example/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

---

## 📦 Services

| Service       | RPCs |
|---------------|------|
| `UserService` | `CreateUser` · `GetUser` · `UpdateUser` · `DeleteUser` · `ListUsers` |
| `PostService` | `CreatePost` · `GetPost` · `UpdatePost` · `DeletePost` · `ListPosts` |

---

## ✅ Validation Rules

All request fields carry [`protovalidate`](https://github.com/bufbuild/protovalidate) constraints encoded directly in the `.proto` files:

| Field | Rule |
|-------|------|
| `id` fields | UUID v4 format |
| `name` | 2–100 chars · alphanumeric, spaces, `_`, `-` |
| `email` | RFC 5322 email format |
| `role` / `status` | Must be a defined enum value (`UNSPECIFIED` disallowed on create) |
| `page` | `>= 1` |
| `page_size` | `1 – 100` |
| `tags` | Max 20 items · each 1–50 chars |
| `content` | Max 100,000 chars |
| Update/filter fields | `IGNORE_IF_ZERO_VALUE` — constraint skipped when field is empty/zero |

---

## 📁 Project Structure

```
proto-package-example/
├── protos/
│   ├── users/v1/
│   │   └── users.proto              # package users.v1 — UserService + validation
│   └── posts/v1/
│       └── posts.proto              # package posts.v1 — PostService + validation
│
├── gen/                             # ⚠️ Generated — do not edit manually
│   ├── ts/                          # TypeScript source (compiled → dist/)
│   │   ├── users/v1/
│   │   │   ├── users_pb.ts          # User messages & enums
│   │   │   └── users_connect.ts     # UserService descriptor
│   │   ├── posts/v1/
│   │   │   ├── posts_pb.ts          # Post messages & enums
│   │   │   └── posts_connect.ts     # PostService descriptor
│   │   └── buf/validate/
│   │       └── validate_pb.ts       # protovalidate runtime types
│   └── go/                          # Go packages (go get-able)
│       ├── users/v1/
│       │   ├── users.pb.go          # package usersv1  — proto messages
│       │   └── usersv1connect/
│       │       └── users.connect.go # package usersv1connect — service interface & client
│       ├── posts/v1/
│       │   ├── posts.pb.go          # package postsv1  — proto messages
│       │   └── postsv1connect/
│       │       └── posts.connect.go # package postsv1connect — service interface & client
│       └── buf/validate/
│           └── validate.pb.go       # protovalidate Go types
│
├── dist/                            # Built JS + .d.ts (published to npm)
│   ├── users/v1/
│   │   ├── users_pb.js / .d.ts
│   │   └── users_connect.js / .d.ts
│   └── posts/v1/
│       ├── posts_pb.js / .d.ts
│       └── posts_connect.js / .d.ts
│
├── buf.yaml                         # Buf workspace config (lint: STANDARD)
├── buf.gen.yaml                     # Code generation plugins
├── buf.lock                         # Locked buf dependency versions
├── go.mod                           # Go module (github.com/sa3akash/proto-package-example)
├── tsconfig.json                    # Compiles gen/ts → dist
└── package.json                     # NPM package (@sa3akash/proto)
```

> **Why `usersv1connect/` and `postsv1connect/`?**
> `protoc-gen-connect-go` intentionally generates the ConnectRPC **service interface, handler, and client** into a separate Go package (`usersv1connect`) from the **proto messages** (`usersv1`). This prevents circular imports — your server implementations import messages from `usersv1` and the service contract from `usersv1connect`.

---

## 🟦 TypeScript Usage

### Install

```bash
npm install @sa3akash/proto @bufbuild/protobuf @connectrpc/connect
```

### Import paths

| Import | Contents |
|--------|----------|
| `@sa3akash/proto/users` | `User`, `UserRole`, all request/response types |
| `@sa3akash/proto/users/connect` | `UserService` ConnectRPC descriptor |
| `@sa3akash/proto/posts` | `Post`, `PostStatus`, all request/response types |
| `@sa3akash/proto/posts/connect` | `PostService` ConnectRPC descriptor |

### Use generated types

```typescript
import { User, UserRole } from "@sa3akash/proto/users";
import { Post, PostStatus } from "@sa3akash/proto/posts";

const user: User = {
  id: "550e8400-e29b-41d4-a716-446655440000",
  name: "Jane Doe",
  email: "jane@example.com",
  role: UserRole.USER,
};
```

### ConnectRPC client (fetch / Node.js)

```typescript
import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-node";
import { UserService } from "@sa3akash/proto/users/connect";
import type { CreateUserRequest } from "@sa3akash/proto/users";
import { UserRole } from "@sa3akash/proto/users";

const transport = createConnectTransport({
  baseUrl: "https://api.example.com",
  httpVersion: "2",
});

const client = createClient(UserService, transport);

const res = await client.createUser({
  name: "Jane Doe",
  email: "jane@example.com",
  role: UserRole.USER,
});

console.log(res.user);
```

### ConnectRPC server handler (Node.js)

```typescript
import { ConnectRouter } from "@connectrpc/connect";
import { UserService } from "@sa3akash/proto/users/connect";
import type { CreateUserRequest, CreateUserResponse } from "@sa3akash/proto/users";

export function registerRoutes(router: ConnectRouter) {
  router.service(UserService, {
    async createUser(req): Promise<CreateUserResponse> {
      return {
        user: {
          id: crypto.randomUUID(),
          name: req.msg.name,
          email: req.msg.email,
          role: req.msg.role,
        },
      };
    },
  });
}
```

### TypeScript validation (optional)

```bash
npm install @bufbuild/protovalidate
```

```typescript
import { createValidator } from "@bufbuild/protovalidate";
import { CreateUserRequest } from "@sa3akash/proto/users";

const validator = await createValidator();

try {
  validator.validate(new CreateUserRequest({ name: "", email: "bad" }));
} catch (err) {
  console.error(err.violations); // field-level violation details
}
```

---

## 🟩 Go Usage

### Install

```bash
go get github.com/sa3akash/proto-package-example
```

### Package layout

| Import path | Package | Contents |
|-------------|---------|----------|
| `.../gen/go/users/v1` | `usersv1` | Proto messages: `User`, `CreateUserRequest`, … |
| `.../gen/go/users/v1/usersv1connect` | `usersv1connect` | `UserServiceClient`, `UserServiceHandler`, `NewUserServiceHandler` |
| `.../gen/go/posts/v1` | `postsv1` | Proto messages: `Post`, `CreatePostRequest`, … |
| `.../gen/go/posts/v1/postsv1connect` | `postsv1connect` | `PostServiceClient`, `PostServiceHandler`, `NewPostServiceHandler` |

### ConnectRPC server

```go
package main

import (
    "context"
    "net/http"

    "connectrpc.com/connect"
    "buf.build/go/protovalidate"

    proto "github.com/sa3akash/proto-package-example/gen/go"
    "github.com/sa3akash/proto-package-example/gen/go/usersv1connect"
)

type UserServer struct{}

func (s *UserServer) CreateUser(
    ctx context.Context,
    req *connect.Request[proto.CreateUserRequest],
) (*connect.Response[proto.CreateUserResponse], error) {
    // Validate request fields against protovalidate constraints
    v, _ := protovalidate.New()
    if err := v.Validate(req.Msg); err != nil {
        return nil, connect.NewError(connect.CodeInvalidArgument, err)
    }

    user := &proto.User{
        Id:    "generated-uuid",
        Name:  req.Msg.Name,
        Email: req.Msg.Email,
        Role:  req.Msg.Role,
    }
    return connect.NewResponse(&proto.CreateUserResponse{User: user}), nil
}

func main() {
    mux := http.NewServeMux()
    path, handler := usersv1connect.NewUserServiceHandler(&UserServer{})
    mux.Handle(path, handler)
    http.ListenAndServe(":8080", mux)
}
```

### ConnectRPC client

```go
import (
    "net/http"
    "connectrpc.com/connect"
    proto "github.com/sa3akash/proto-package-example/gen/go/users/v1"
    "github.com/sa3akash/proto-package-example/gen/go/users/v1/usersv1connect"
)

client := usersv1connect.NewUserServiceClient(
    http.DefaultClient,
    "https://api.example.com",
)

resp, err := client.GetUser(ctx, connect.NewRequest(&proto.GetUserRequest{
    Id: "550e8400-e29b-41d4-a716-446655440000",
}))
if err != nil {
    log.Fatal(err)
}
fmt.Println(resp.Msg.User)
```

---

## 🔧 Local Development

### Prerequisites

- [Buf CLI](https://buf.build/product/cli) — installed via npm: `npx @bufbuild/buf`  
- [Bun](https://bun.sh) or Node.js 20+
- [Go 1.22+](https://golang.org/dl/)

### Generate TypeScript + Go code

```bash
# First time — fetch buf dependencies
bun run generate   # runs: buf generate
# or directly:
bunx buf dep update
bunx buf generate
```

Generated files appear in `gen/ts/` and `gen/go/`.

### Build TypeScript package

```bash
bun install
bun run build      # runs: tsc -p tsconfig.json
```

### Verify Go module compiles

```bash
go mod tidy
go build ./...
```

---

## 🚀 Publishing to NPM

Push a version tag — GitHub Actions handles the rest:

```bash
git tag v1.0.1
git push origin v1.0.1
```

The `publish.yml` workflow will:
1. Run `buf generate` (fresh TS + Go code)
2. Build TypeScript (`tsc`)
3. Stamp version from the git tag
4. Publish `@sa3akash/proto@1.0.1` to **npmjs.com**

### Required GitHub Secret

| Secret | Description |
|--------|-------------|
| `NPM_TOKEN` | npm automation token with publish access |

> Create at: <https://www.npmjs.com/settings/tokens>

# Refactoring and Dependency Injection Plan

This document outlines the plan to refactor the Peril codebase for improved quality, structure, and to introduce the Google Wire dependency injection container.

## Phase 1: Code Restructuring and Service Encapsulation

### Step 1.1: Analyze Existing Codebase
- **Objective:** Understand the current dependency flow and identify areas for refactoring.
- **Action:** Use `codebase_investigator` to map out the current structure. Pay close attention to `cmd/server/main.go` and `cmd/client/main.go` to see how dependencies like `pubsub`, `gamelogic`, and `routing` are being created and used.

### Step 1.2: Create a Dedicated `users` Package
- **Objective:** Encapsulate user-related logic into its own package to improve separation of concerns.
- **Action:**
    1. Create a new directory: `internal/users`.
    2. Create a new file: `internal/users/users.go`.
    3. Move the `User` struct and any related functions from `internal/routing/models.go` to `internal/users/users.go`.
    4. Update import paths in files that use the `User` struct.

### Step 1.3: Create a `boot` Package for Server Initialization
- **Objective:** Centralize the server startup and dependency wiring logic, cleaning up the `cmd/server/main.go` file.
- **Action:**
    1. Create a new directory: `internal/boot`.
    2. Create a new file: `internal/boot/server.go`.
    3. In `internal/boot/server.go`, create a function `Startup` that will contain the logic for initializing the router, pub/sub connection, and other server-side components. This will involve moving code from `cmd/server/main.go`.

### Step 1.4: Refactor Server's Main Function
- **Objective:** Simplify the server's entry point to only call the new startup logic.
- **Action:**
    1. Modify `cmd/server/main.go` to call the `boot.Startup()` function. The `main` function should become very minimal.

### Step 1.5: Create a `boot` Package for Client Initialization
- **Objective:** Centralize the client startup and dependency wiring logic, cleaning up the `cmd/client/main.go` file.
- **Action:**
    1. Create a new file: `internal/boot/client.go`.
    2. In `internal/boot/client.go`, create a function `ClientStartup` that will handle the client's initialization logic, moved from `cmd/client/main.go`.

### Step 1.6: Refactor Client's Main Function
- **Objective:** Simplify the client's entry point.
- **Action:**
    1. Modify `cmd/client/main.go` to call the `boot.ClientStartup()` function.

## Phase 2: Introducing Dependency Injection with Wire

### Step 2.1: Install Wire
- **Objective:** Add the Wire tool to the project.
- **Action:** Run the command `go get github.com/google/wire/cmd/wire`.

### Step 2.2: Create Providers for Services
- **Objective:** Create provider functions for each of the core components.
- **Action:**
    1. In `internal/pubsub/pubsub.go`, create a `NewPubSub` function that returns a `*PubSub`.
    2. In `internal/gamelogic/gamelogic.go`, create a `NewGameLogic` function.
    3. In `internal/routing/routing.go`, create a `NewRouter` function.
    4. ... and so on for other components.

### Step 2.3: Implement Wire on the Server-Side
- **Objective:** Use Wire to manage dependencies in the server.
- **Action:**
    1. Create a file `cmd/server/wire.go`.
    2. In this file, define the `ProviderSet` and the `InitializeServer` injector function.
    3. Run `wire` in the `cmd/server` directory.
    4. Update `internal/boot/server.go` to use the generated `InitializeServer` function.

### Step 2.4: Implement Wire on the Client-Side
- **Objective:** Use Wire to manage dependencies in the client.
- **Action:**
    1. Create a file `cmd/client/wire.go`.
    2. Define the `ProviderSet` and `InitializeClient` injector function.
    3. Run `wire` in the `cmd/client` directory.
    4. Update `internal/boot/client.go` to use the generated `InitializeClient` function.

### Step 2.5: Final Validation
- **Objective:** Ensure the application functions as before.
- **Action:**
    1. Run the server and client.
    2. Verify they can communicate and that the game logic is working.
    3. Run `go mod tidy` to clean up dependencies.

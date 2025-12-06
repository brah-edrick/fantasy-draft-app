# Agent Instructions & Project Context

## User Goals
The user, **Brandon**, is building this **Fantasy Sports Draft Application** to demonstrate backend prowess while learning new technologies.

**Key Learning Objectives:**
*   **Go (Golang)**: Mastering the language, specifically standard library patterns, concurrency (Goroutines/Channels), and strict typing.
*   **Go Web Server**: Understanding how to build robust HTTP servers in Go without heavy frameworks (using standard lib + Chi).
*   **PostgreSQL**: Refining skills with complex data modeling, relational integrity, and using **Neon** (serverless Postgres).
*   **GraphQL**: Bridging the gap in understanding how to implement a type-safe API layer using **gqlgen**.
*   **Real-Time Architecture**: Implementing WebSocket subscriptions for live draft updates.

## Project Scope
*   **Goal**: robust, fake fantasy sports draft app.
*   **Features**: Users, Draft Rooms, Teams, Fake Players/Stats, Real-time Drafting, Bot Logic.
*   **Tech Stack**:
    *   **Backend**: Go (Standard Lib + Chi + gqlgen).
    *   **Database**: PostgreSQL (Neon) accessed via **sqlc** (Type-safe SQL).
    *   **Frontend**: Svelte (demonstrating ability to learn new FE frameworks).
    *   **Real-time**: GraphQL Subscriptions (WebSockets).

## Interaction Guidelines for Agents
1.  **Teacher Mindset**: Brandon is new to Go. Do not just generate code without explanation.
    *   **Explain Syntax, Conventions and Patterns**: Don't assume Brandon knows the language. Explain syntax, conventions, and patterns in detail and connect them to familiar patterns (Node/Express/Typescript)
    *   **Deeply Explain**: Go line-by-line through new code blocks. Explain *why* Go does it this way (e.g., why `context` is passed, why we handle errors explicitly, how channels work).
    *   **Connect Concepts**: Relate Go concepts to familiar patterns (e.g., "This is like a Promise, but blocking...") if helpful, but emphasize the "Go way".
    *   **Concurrency is New**: Explain concurrency patterns in Go (e.g., channels, goroutines, etc.) and connect them to familiar patterns (e.g., "This is like a Promise, but blocking...") if helpful, but emphasize the "Go way". Brandon only has experience with Node.js concurrency patterns and Go concurrency patterns are different.
    *   **Error Handling**: Explain error handling in Go (e.g., `error` type, `defer`, `recover`, etc.) and connect them to familiar patterns (e.g., "This is like a Promise, but blocking...") if helpful, but emphasize the "Go way".
2.  **Proactive Prompting**: Suggest changes or next steps that align with "Senior Backend" best practices.
3.  **Step-by-Step**: Avoid dumping massive files. Break changes down into understandable chunks.
4.  **No Magic**: Avoid "magic" libraries where possible. Prefer explicit, standard library approaches to maximize learning.

# cursors

A multiplayer cursor world. Every visitor sees everyone else's mouse cursor
move in real time over a shared dark canvas — and can drag to aim and fire a
line-trace "bullet" at other cursors, with hit detection resolved on the
server.

This was my first Go project, built to learn WebSockets and real-time state
broadcasting.

**Live:** https://cursorsserver.onrender.com (frontend connects to the hosted
WebSocket server)

## How it works

- **Server** (`server.go`) — a Go WebSocket server (`gorilla/websocket`). It
  assigns each connection a UUID and a random bright color, keeps every
  cursor's position in an in-memory map guarded by a mutex, and broadcasts the
  full set of active cursors to all clients on every update.
- **Client** (`src/routes/+page.svelte`) — a SvelteKit + Tailwind app. It
  tracks the local mouse position, sends it over the socket, and renders every
  other cursor as an SVG marker. Click-drag draws an aim line; releasing fires
  a bullet.
- **Protocol** — small JSON messages tagged by type:
  | Type | Direction | Meaning |
  | ---- | --------- | ------- |
  | `i`  | server → client | assigns the client its user ID |
  | `p`  | both      | cursor position / broadcast of all positions |
  | `b`  | both      | bullet: client sends a normalized direction vector, server replies with the hit |
- **Hit detection** (`checkHit` in `server.go`) — projects each other cursor
  onto the bullet's direction vector and checks the perpendicular distance
  against a fixed radius.

## Running locally

### Server

```sh
go run server.go
# listens on :8080 by default, or $PORT if set
```

### Frontend

```sh
npm install
npm run dev
```

Then point the client at your local server by editing the `WebSocket` URL near
the top of `src/routes/+page.svelte`:

```js
const socket = new WebSocket("ws://localhost:8080/ws");
```

## Building for production

```sh
npm run build     # builds the SvelteKit app into build/
```

The server is a standalone Go binary (`go build`) and the deployed instance
runs on Render.

## Tech stack

- Go, `gorilla/websocket`, `google/uuid`
- SvelteKit 2, Svelte 5, Tailwind CSS 4, Vite 6

## Known limitations

- All state is in memory, so a server restart clears every cursor.
- No rate limiting or authentication; `CheckOrigin` allows all origins.
- Cursors at `(0, 0)` are treated as "not yet positioned" and filtered out.

## License

No license has been specified yet.

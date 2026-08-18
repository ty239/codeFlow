# codeFlow

**codeFlow** is a visual whiteboard for data structures and algorithms — think Excalidraw, but purpose-built for drawing and animating the way DSA problems (like the ones on LeetCode) actually work.

## The problem

Whiteboarding a LeetCode problem — a linked list reversal, a graph traversal, a DP table filling in — helps me understand it far better than reading code line by line. General drawing tools like Excalidraw work, but they're not built for this: there's no concept of a "node," no way to step through an algorithm frame by frame, and no primitives for the shapes DSA problems actually need (arrays, trees, graphs, pointers, stacks).

## The goal

Build a tool that's:
- **Purpose-made for DSA** — first-class shapes for arrays, linked lists, trees, graphs, and stacks/queues, instead of generic rectangles and arrows.
- **Flow-aware** — able to show step-by-step execution (e.g. two pointers moving, a recursion tree unfolding, a BFS frontier expanding), not just a static diagram.
- **More accessible than Excalidraw** — lower friction for this specific use case: fewer clicks to drop in a data structure, sensible defaults, less manual drawing.
- **A personal study tool first** — I'm building this to visualize problems while I practice, with the side goal of it being useful to others learning DSA the same way.

## Status

Early stage. Backend has a Postgres-backed schema (users, boards, elements, connections, collaborators) and user signup/login (bcrypt + JWT). Boards/elements CRUD, the canvas, and data-structure primitives don't exist yet. Treat this README as a statement of intent/roadmap rather than a description of finished features.

## Planned features

- [ ] Canvas/whiteboard surface (drawing, shapes, freehand — Excalidraw-style)
- [ ] Built-in DSA primitives: arrays, linked lists, binary trees, graphs, stacks, queues
- [ ] Step-through / playback mode to animate an algorithm's execution over a data structure
- [ ] Save/load boards (per problem, so you can build a personal library of visualized solutions)
- [ ] Export (image/link) to share a visualization

## Tech stack

- **Backend:** Go (`net/http`, `pgx`/Postgres, `bcrypt`, `golang-jwt`)
- **Frontend:** TBD

## Why this project

Built to solve my own problem while practicing LeetCode: visualizing algorithms and data structures is how I actually learn them, and I wanted a tool tailored to that instead of fighting a general-purpose drawing app.

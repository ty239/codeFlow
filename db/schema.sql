-- =====================================================================
-- codeFlow database schema
-- =====================================================================
-- A collaborative whiteboard for visualizing data structures/algorithms.
--
-- Table overview:
--   users         - registered accounts
--   boards        - a whiteboard/canvas, owned by one user
--   elements      - the shapes drawn on a board (nodes, arrays, labels, ...)
--   connections   - arrows/edges linking two elements (pointers, edges, ...)
--   collaborators - join table granting other users access to a board
--
-- Relationships:
--   users (1) --- (N) boards            [boards.owner_id]
--   boards (1) --- (N) elements         [elements.board_id]
--   elements (1) --- (N) elements       [elements.parent_id, for grouping]
--   boards (1) --- (N) connections      [connections.board_id]
--   elements (1) --- (N) connections    [connections.source_id / target_id]
--   boards (N) --- (N) users            [collaborators]
-- =====================================================================


-- Registered user accounts.
CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username      VARCHAR(50)  NOT NULL UNIQUE,
    email         VARCHAR(255) NOT NULL UNIQUE,
    name          VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now()  -- when the account was created
);


-- A whiteboard/canvas. Each board belongs to exactly one owner and can
-- additionally be shared with other users via the `collaborators` table.
CREATE TABLE boards (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),  -- random UUID so board URLs/ids aren't guessable
    owner_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- creator of the board; delete their boards if the user is deleted
    title      VARCHAR(255) NOT NULL DEFAULT 'Untitled',  -- display name of the board
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),   -- when the board was created
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()    -- when the board was last modified
);


-- Core visual element on a board (nodes, shapes, array cells, labels, etc.)
-- Elements can be nested via parent_id (e.g. array cells belong to an
-- array container), which lets complex structures be grouped and moved
-- together.
CREATE TABLE elements (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    board_id   UUID NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    parent_id  UUID REFERENCES elements(id) ON DELETE CASCADE,  -- for grouping (array cells belong to an array container)

    -- what kind of element
    type       VARCHAR(50) NOT NULL,  -- 'node', 'array', 'cell', 'label', 'stack', 'grid', etc.

    -- visual / spatial (position and size on the canvas)
    x          DOUBLE PRECISION NOT NULL DEFAULT 0,
    y          DOUBLE PRECISION NOT NULL DEFAULT 0,
    width      DOUBLE PRECISION NOT NULL DEFAULT 100,
    height     DOUBLE PRECISION NOT NULL DEFAULT 100,
    rotation   DOUBLE PRECISION NOT NULL DEFAULT 0,
    z_index    INTEGER NOT NULL DEFAULT 0,  -- stacking order; higher draws on top

    -- content
    value      TEXT,           -- what's displayed inside ("5", "null", "head", etc.)
    label      TEXT,           -- optional annotation ("i", "left", "root")

    -- styling
    shape      VARCHAR(30) DEFAULT 'rectangle',  -- 'circle', 'rectangle', 'diamond', 'roundedRect'
    fill       VARCHAR(30) DEFAULT '#ffffff',    -- fill color (hex)
    stroke     VARCHAR(30) DEFAULT '#000000',    -- border color (hex)
    font_size  INTEGER DEFAULT 16,

    -- type-specific stuff (array order, grid row/col, etc.), stored loosely
    -- so new element types don't require schema changes
    properties JSONB NOT NULL DEFAULT '{}',

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Arrows/edges connecting two elements (linked list pointers, tree edges, graph edges)
CREATE TABLE connections (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    board_id   UUID NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    source_id  UUID NOT NULL REFERENCES elements(id) ON DELETE CASCADE,  -- element the arrow starts at
    target_id  UUID NOT NULL REFERENCES elements(id) ON DELETE CASCADE,  -- element the arrow points to

    type       VARCHAR(30) NOT NULL DEFAULT 'directed',  -- 'directed', 'undirected', 'double'
    label      TEXT,          -- edge weight, pointer name ("next", "left", "right")
    style      VARCHAR(20) DEFAULT 'solid',  -- 'solid', 'dashed', 'dotted'
    stroke     VARCHAR(30) DEFAULT '#000000',  -- line color (hex)

    properties JSONB NOT NULL DEFAULT '{}',  -- type-specific extras, stored loosely like elements.properties

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Who can access/edit a board besides the owner. Enables multiple users
-- editing (or just viewing) the same board at the same time.
CREATE TABLE collaborators (
    board_id   UUID NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role       VARCHAR(20) NOT NULL DEFAULT 'editor',  -- 'editor', 'viewer'
    joined_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (board_id, user_id)  -- one entry per user per board
);

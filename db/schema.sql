CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username      VARCHAR(50)  NOT NULL UNIQUE,
    email         VARCHAR(255) NOT NULL UNIQUE,
    name          VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now()  --The time of when the users does somthin
);


CREATE TABLE boards (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),  -- genrtest radom sequnces of unique url/id so non one can acess it
    owner_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- forgien key make the sure the data belongs to the user and if session is deleted delte
    title      VARCHAR(255) NOT NULL DEFAULT 'Untitled',  -- name of the borad
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),   -- time created 
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()   -- time updated
);



                  
-- Core visual element on a board (nodes, shapes, array cells, labels, etc.)
CREATE TABLE elements (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    board_id   UUID NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    parent_id  UUID REFERENCES elements(id) ON DELETE CASCADE,  -- for grouping (array cells belong to an array container)
    
    -- what kind of element
    type       VARCHAR(50) NOT NULL,  -- 'node', 'array', 'cell', 'label', 'stack', 'grid', etc.
    
    -- visual / spatial
    x          DOUBLE PRECISION NOT NULL DEFAULT 0,
    y          DOUBLE PRECISION NOT NULL DEFAULT 0,
    width      DOUBLE PRECISION NOT NULL DEFAULT 100,
    height     DOUBLE PRECISION NOT NULL DEFAULT 100,
    rotation   DOUBLE PRECISION NOT NULL DEFAULT 0,
    z_index    INTEGER NOT NULL DEFAULT 0,
    
    -- content
    value      TEXT,           -- what's displayed inside ("5", "null", "head", etc.)
    label      TEXT,           -- optional annotation ("i", "left", "root")
    
    -- styling
    shape      VARCHAR(30) DEFAULT 'rectangle',  -- 'circle', 'rectangle', 'diamond', 'roundedRect'
    fill       VARCHAR(30) DEFAULT '#ffffff',
    stroke     VARCHAR(30) DEFAULT '#000000',
    font_size  INTEGER DEFAULT 16,
    
    -- type-specific stuff (array order, grid row/col, etc.)
    properties JSONB NOT NULL DEFAULT '{}',
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Arrows/edges connecting two elements (linked list pointers, tree edges, graph edges)
CREATE TABLE connections (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    board_id   UUID NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    source_id  UUID NOT NULL REFERENCES elements(id) ON DELETE CASCADE,
    target_id  UUID NOT NULL REFERENCES elements(id) ON DELETE CASCADE,
    
    type       VARCHAR(30) NOT NULL DEFAULT 'directed',  -- 'directed', 'undirected', 'double'
    label      TEXT,          -- edge weight, pointer name ("next", "left", "right")
    style      VARCHAR(20) DEFAULT 'solid',  -- 'solid', 'dashed', 'dotted'
    stroke     VARCHAR(30) DEFAULT '#000000',
    
    properties JSONB NOT NULL DEFAULT '{}',
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Who can access/edit a board besides the owner
CREATE TABLE collaborators (
    board_id   UUID NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role       VARCHAR(20) NOT NULL DEFAULT 'editor',  -- 'editor', 'viewer'
    joined_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    
    PRIMARY KEY (board_id, user_id)  -- one entry per user per board
);                -- elements shapes go deeper using tree array shapes 
                  -- collaborators 2 people can draw at the exact same time
 
                        




























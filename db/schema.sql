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



                  -- drawing table
                  -- elements shapes go deeper using tree array shapes 
                  -- collaborators 2 people can draw at the exact same time
 
                        




























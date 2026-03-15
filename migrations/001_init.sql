CREATE TABLE IF NOT EXISTS black_list (
    id SERIAL PRIMARY KEY,
    subnet CIDR NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS white_list (
    id SERIAL PRIMARY KEY,
    subnet CIDR NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_black_list_subnet ON black_list USING gist (subnet inet_ops);
CREATE INDEX idx_white_list_subnet ON white_list USING gist (subnet inet_ops);

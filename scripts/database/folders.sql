CREATE TABLE folders (
    id SERIAL,
    parent_id INT,
    name VARCHAR(60) NOT NULL,
    created_at TIMESTAMP default CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL,
    deleted BOOLEAN NOT NULL DEFAULT FALSE
    PRIMARY KEY (id)
    CONSTRAINT fk_folders_parent_id FOREIGN KEY (parent_id) REFERENCES folders(id)
);

CREATE TABLE files (
    id SERIAL,
    folder_id INT,
    owner_id INT NOT NULL,
    name VARCHAR(200) NOT NULL,
    type VARCHAR(50) NOT NULL,
    path VARCHAR(255) NOT NULL,
    created_at TIMESTAMP default CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL,
    deleted BOOLEAN NOT NULL DEFAULT FALSE
    PRIMARY KEY (id)
    CONSTRAINT fk_files_folder_id FOREIGN KEY (folder_id) REFERENCES folders(id),
    CONSTRAINT fk_files_owner_id FOREIGN KEY (owner_id) REFERENCES users(id)
);

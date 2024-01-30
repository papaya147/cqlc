CREATE TABLE data.papaya(
    id int,
    name varchar NOT NULL,
    price int NOT NULL,
    description TEXT,
    image_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (
        (id),
        name,
        price,
        description,
        image_url,
        created_at,
        updated_at
    )
);
CREATE TABLE IF NOT EXISTS quota (
        id SERIAL,
	quota INTEGER, 
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
)

package migrations

type Tags struct{}

func (*Tags) up() []string {
	return []string{
		`CREATE TABLE IF NOT EXISTS tags (
            id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
            name VARCHAR(64) NOT NULL UNIQUE,
            -- user_id BIGINT NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc')
            -- FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
        );`, `
        CREATE OR REPLACE FUNCTION getTagIdByName(tagName VARCHAR(16)) RETURNS BIGINT 
        AS $$ 
        DECLARE i BIGINT;
        BEGIN 
            SELECT id INTO i FROM tags WHERE name=tagName;
            IF i IS NULL THEN 
                RAISE EXCEPTION 'tag not found'; 
            END IF; 
            RETURN i;
        END; 
        $$ LANGUAGE PLPGSQL;`,
	}
}

func (*Tags) down() []string {
	return []string{
		`DROP FUNCTION IF EXISTS getTagIdByName;`,
		`DROP TABLE IF EXISTS tags;`,
	}
}

func (*Tags) tableName() string {
	return "tags"
}
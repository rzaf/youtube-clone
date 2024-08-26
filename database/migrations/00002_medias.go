package migrations

type Medias struct{}

func (*Medias) up() []string {
	return []string{
		`
        CREATE TABLE IF NOT EXISTS medias (
            id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
            title VARCHAR(128) NOT NULL,
            text VARCHAR(1024) NULL,
            url VARCHAR(16) NOT NULL UNIQUE,
            media_type SMALLINT NOT NULL, --0:video 1:music 2:photo
            thumbnail VARCHAR(16) NULL,
            user_id BIGINT NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
            updated_at TIMESTAMP NULL,
            FOREIGN KEY ("user_id") REFERENCES "users" ("id")
        );
        COMMENT ON COLUMN medias.media_type IS '0:video 1:music 2:photo';
        `, `
        CREATE OR REPLACE FUNCTION getMediaIdByUrl(mediaUrl VARCHAR(16)) RETURNS BIGINT 
        AS $$ 
        DECLARE i BIGINT;
        BEGIN 
            SELECT id INTO i FROM medias WHERE url=mediaUrl;
            IF i IS NULL THEN 
                RAISE EXCEPTION 'media not found'; 
            END IF; 
            RETURN i;
        END; 
        $$ LANGUAGE PLPGSQL;
        `, `
        CREATE OR REPLACE FUNCTION getMediaIdByUrlAndType(mediaUrl VARCHAR(16),mediaType INT) RETURNS BIGINT 
        AS $$ 
        DECLARE i BIGINT;
        DECLARE t INT;
        BEGIN 
            SELECT id,media_type INTO i,t FROM medias WHERE url=mediaUrl;
            IF i IS NULL THEN 
                RAISE EXCEPTION 'media not found'; 
            END IF; 
            IF mediatype!=3 AND t != mediaType THEN 
                RAISE EXCEPTION 'ivalid media type'; 
            END IF; 
            RETURN i;
        END; 
        $$ LANGUAGE PLPGSQL;
        `,
	}
}

func (*Medias) down() []string {
	return []string{
		`DROP FUNCTION IF EXISTS getMediaIdByUrl;`,
		`DROP TABLE IF EXISTS medias;`,
	}
}

func (*Medias) tableName() string {
	return "medias"
}
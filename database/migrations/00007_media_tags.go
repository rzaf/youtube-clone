package migrations

type MediaTags struct{}

func (*MediaTags) up() []string {
	return []string{
		`CREATE TABLE IF NOT EXISTS media_tags (
            id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
            media_id BIGINT NOT NULL,
            tag_id BIGINT NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
            FOREIGN KEY (media_id) REFERENCES medias (id) ON DELETE CASCADE,
            FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE,
            UNIQUE(media_id,tag_id)
        );`,
	}
}

func (*MediaTags) down() []string {
	return []string{
		`DROP TABLE IF EXISTS media_tags;`,
	}
}

func (*MediaTags) tableName() string {
	return "media_tags"
}
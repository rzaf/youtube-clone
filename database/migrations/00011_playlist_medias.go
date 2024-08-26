package migrations

type PlaylistMedias struct{}

func (*PlaylistMedias) up() []string {
	return []string{
		`CREATE TABLE IF NOT EXISTS playlist_medias (
            id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
            playlist_id BIGINT NOT NULL,
            text VARCHAR(512) NULL,
            media_id BIGINT NOT NULL,
            custom_order INT NOT NULL DEFAULT 1,
            created_at TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
            updated_at TIMESTAMP NULL,
            UNIQUE(playlist_id,media_id),
            FOREIGN KEY (media_id) REFERENCES medias (id) ON DELETE CASCADE,
            FOREIGN KEY (playlist_id) REFERENCES playlists (id) ON DELETE CASCADE
        );
        COMMENT ON COLUMN playlist_medias.text IS'playlister additional note on video';`,
	}
}

func (*PlaylistMedias) down() []string {
	return []string{
		`DROP TABLE IF EXISTS playlist_medias;`,
	}
}

func (*PlaylistMedias) tableName() string {
	return "playlist_medias"
}

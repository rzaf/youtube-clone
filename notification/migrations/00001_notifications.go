package migrations

// up

type notifications struct{}

func (*notifications) up() []string {
	return []string{
		`
        CREATE TABLE IF NOT EXISTS notifications (
            id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
            user_id BIGINT NOT NULL,
            title VARCHAR(225) NOT NULL,
            message VARCHAR(1000) NOT NULL,
            -- data json,
            seen_at TIMESTAMP NULL,
            created_at TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
            updated_at TIMESTAMP NULL
        );`,
	}
}

func (*notifications) down() []string {
	return []string{
		`DROP TABLE IF EXISTS notifications;`,
	}
}

func (*notifications) tableName() string {
	return "notifications"
}

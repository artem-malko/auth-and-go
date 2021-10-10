package repository

import migrate "github.com/rubenv/sql-migrate"

// Migrations a list of migrations to create/update/rollback database
var Migrations = []*migrate.Migration{
	{
		Id: "20200224000001_create_session_client_id_type",
		Up: []string{
			`
			CREATE TYPE client_id_type AS ENUM ('web', 'native_android', 'native_ios')`,
		},
		Down: []string{
			`
			DROP TYPE client_id_type CASCADE;
			`,
		},
	},
	{
		Id: "20200224000002_create_sessions_table",
		Up: []string{
			`
			CREATE TABLE sessions
			(
				id 							uuid not null primary key,
				account_id 					uuid not null,
				identity_id					uuid not null,
				client_id					client_id_type not null,
				access_token 				uuid not null,
				access_token_expires_date   timestamp with time zone not null,
				refresh_token 				uuid not null,
				refresh_token_expires_date  timestamp with time zone not null,

				CONSTRAINT account_id_fkey FOREIGN KEY (account_id)
					REFERENCES accounts (id) MATCH FULL
				 	ON DELETE CASCADE,

				CONSTRAINT identity_id_fkey FOREIGN KEY (identity_id)
					REFERENCES identities (id) MATCH FULL
				 	ON DELETE CASCADE
			);
			`,
		},
	},
	{
		Id: "20200224000003_create_sessions_table_refresh_token_index",
		Up: []string{
			`
			CREATE UNIQUE INDEX
				sessions_refresh_token_uniq_idx ON sessions (refresh_token);
			`,
		},
		Down: []string{
			`
			DROP INDEX sessions_refresh_token_uniq_idx;
			`,
		},
	},
	{
		Id: "20200224000004_create_sessions_table_access_token_index",
		Up: []string{
			`
			CREATE UNIQUE INDEX
				sessions_access_token_uniq_idx ON sessions (access_token);
			`,
		},
		Down: []string{
			`
			DROP INDEX sessions_access_token_uniq_idx;
			`,
		},
	},
	{
		Id: "20200224000005_create_sessions_table_account_id_index",
		Up: []string{
			`
			CREATE INDEX
				sessions_account_id_idx ON sessions USING HASH (account_id);
			`,
		},
		Down: []string{
			`
			DROP INDEX sessions_account_id_idx;
			`,
		},
	},
	{
		Id: "20200224000006_create_sessions_table_refresh_token_expires_date_index",
		Up: []string{
			`
			CREATE INDEX
				sessions_refresh_token_expires_date_idx ON sessions ((refresh_token_expires_date::timestamp with time zone) DESC);
			`,
		},
		Down: []string{
			`
			DROP INDEX sessions_refresh_token_expires_date_idx;
			`,
		},
	},
}

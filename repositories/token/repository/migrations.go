package repository

import (
	migrate "github.com/rubenv/sql-migrate"
)

// Migrations a list of migrations to create/update/rollback database
var Migrations = []*migrate.Migration{
	{
		Id: "20200424000000_create_token_status_type",
		Up: []string{
			`
			CREATE TYPE token_status AS ENUM ('active', 'used')`,
		},
		Down: []string{
			`
			DROP TYPE token_status CASCADE;
			`,
		},
	},
	{
		Id: "20200424000001_create_token_type",
		Up: []string{
			`
			CREATE TYPE token_type AS ENUM ('registration_confirmation', 'email_confirmation', 'auto_login')`,
		},
		Down: []string{
			`
			DROP TYPE token_type CASCADE;
			`,
		},
	},
	{
		Id: "20200424000002_create_tokens_table",
		Up: []string{
			`
			CREATE TABLE tokens
			(
				id 						uuid not null primary key,
				token_type				token_type not null,
				token_status			token_status not null default 'active',
				account_id 				uuid not null,
				identity_id				uuid not null,
				client_id				client_id_type not null,
				expires_date	 		timestamp with time zone not null default now() + interval '1 hour',

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
		Id: "20200424000003_create_tokens_table_registration_confirmation_token_expires_date_index",
		Up: []string{
			`
			CREATE INDEX
				tokens_registration_confirmation_token_expires_date_idx 
				ON tokens ((expires_date::timestamp with time zone) DESC)
				WHERE token_type = 'registration_confirmation' AND token_status = 'active';
			`,
		},
		Down: []string{
			`
			DROP INDEX tokens_registration_confirmation_token_expires_date_idx;
			`,
		},
	},
	{
		Id: "20200424000004_create_tokens_table_email_confirmation_token_expires_date_index",
		Up: []string{
			`
			CREATE INDEX
				tokens_email_confirmation_token_expires_date_idx 
				ON tokens ((expires_date::timestamp with time zone) DESC)
				WHERE token_type = 'email_confirmation' AND token_status = 'active';
			`,
		},
		Down: []string{
			`
			DROP INDEX tokens_email_confirmation_token_expires_date_idx;
			`,
		},
	},
	{
		Id: "20200424000005_create_tokens_table_auto_login_token_expires_date_index",
		Up: []string{
			`
			CREATE INDEX
				tokens_auto_login_token_expires_date_idx 
				ON tokens ((expires_date::timestamp with time zone) DESC)
				WHERE token_type = 'auto_login' AND token_status = 'active';
			`,
		},
		Down: []string{
			`
			DROP INDEX tokens_auto_login_token_expires_date_idx;
			`,
		},
	},
	{
		Id: "20200424000006_create_tokens_table_used_tokens_index",
		Up: []string{
			`
			CREATE INDEX tokens_used_tokens_idx ON tokens (token_status) WHERE token_status = 'used';
			`,
		},
		Down: []string{
			`
			DROP INDEX tokens_used_tokens_idx;
			`,
		},
	},
	{
		Id: "20200424000007_create_tokens_table_identity_id_index",
		Up: []string{
			`
			CREATE INDEX
				tokens_identity_id_idx ON tokens USING HASH (identity_id);
			`,
		},
		Down: []string{
			`
			DROP INDEX tokens_identity_id_idx;
			`,
		},
	},
}

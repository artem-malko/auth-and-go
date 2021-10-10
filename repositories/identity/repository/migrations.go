package repository

import (
	migrate "github.com/rubenv/sql-migrate"
)

// Migrations a list of migrations to create/update/rollback database
var Migrations = []*migrate.Migration{
	{
		Id: "20200221000001_create_identity_type",
		Up: []string{
			`
			CREATE TYPE identity_type AS ENUM ('email', 'facebook', 'google')`,
		},
		Down: []string{
			`
			DROP TYPE identity_type CASCADE;
			`,
		},
	},
	{
		Id: "20200221000002_create_identity_status_type",
		Up: []string{
			`
			CREATE TYPE identity_status AS ENUM ('unconfirmed', 'confirmed')`,
		},
		Down: []string{
			`
			DROP TYPE identity_status CASCADE;
			`,
		},
	},
	{
		Id: "20200221000003_create_identities_table",
		Up: []string{
			`
			CREATE TABLE  identities
			(
				id 						uuid not null primary key,
				account_id 				uuid not null,
				identity_status			identity_status not null default 'unconfirmed',
				identity_type			identity_type not null,
				google_social_id		varchar(250),
				facebook_social_id		varchar(250),
				email					text,
				password_hash			text,
				created_at	 			timestamp with time zone not null default now(),
				updated_at	 			timestamp with time zone not null default now(),

				CONSTRAINT account_id_fkey FOREIGN KEY (account_id)
					REFERENCES accounts (id) MATCH FULL
				 	ON DELETE CASCADE
			);
			`,
		},
	},
	{
		Id: "20200221000004_create_identities_table_google_social_id_uniq_idx_index",
		Up: []string{
			`
			CREATE UNIQUE INDEX
				identities_google_social_id_uniq_idx ON identities (google_social_id)
			WHERE identity_type = 'google';
			`,
		},
		Down: []string{
			`
			DROP INDEX identities_google_social_id_uniq_idx;
			`,
		},
	},
	{
		Id: "20200221000005_create_identities_table_facebook_social_id_uniq_idx_index",
		Up: []string{
			`
			CREATE UNIQUE INDEX
				identities_facebook_social_id_uniq_idx ON identities (facebook_social_id)
			WHERE identity_type = 'facebook';
			`,
		},
		Down: []string{
			`
			DROP INDEX identities_facebook_social_id_uniq_idx;
			`,
		},
	},
	{
		Id: "20200221000006_create_identities_table_email_uniq_idx_index",
		Up: []string{
			`
			CREATE UNIQUE INDEX
				identities_email_uniq_idx ON identities (email) WHERE identity_type = 'email';
			`,
		},
		Down: []string{
			`
			DROP INDEX identities_email_uniq_idx;
			`,
		},
	},
	{
		Id: "20200221000007_create_identities_table_email_idx_index",
		Up: []string{
			`
			CREATE INDEX
				identities_email_idx ON identities (email);
			`,
		},
		Down: []string{
			`
			DROP INDEX identities_email_idx;
			`,
		},
	},
}

package repository

import (
	migrate "github.com/rubenv/sql-migrate"
)

// Migrations a list of migrations to create/update/rollback database
var Migrations = []*migrate.Migration{
	{
		Id: "20200121000001_create_account_status_type",
		Up: []string{
			`
			CREATE TYPE account_status AS ENUM ('unconfirmed', 'confirmed', 'deleted', 'banned')`,
		},
		Down: []string{
			`
			DROP TYPE account_status CASCADE;
			`,
		},
	},
	{
		Id: "20200121000002_create_account_type",
		Up: []string{
			`
			CREATE TYPE account_type AS ENUM ('free', 'commercial')`,
		},
		Down: []string{
			`
			DROP TYPE account_type CASCADE;
			`,
		},
	},
	{
		Id: "20200121000003_create_accounts_table",
		Up: []string{
			`
			CREATE TABLE accounts
			(
				id 						uuid not null primary key,
				account_type			account_type not null,
				account_status			account_status not null default 'unconfirmed',
				account_name			text not null,
				profile					jsonb,
				settings				jsonb,
				last_ip					text,
				last_login	 			timestamp with time zone not null default now(),
				created_at	 			timestamp with time zone not null default now(),
				updated_at	 			timestamp with time zone not null default now()
			);
			`,
		},
	},
	{
		Id: "20200122000004_create_accounts_table_account_name_active_account_uniq_idx_index",
		Up: []string{
			`
			CREATE UNIQUE INDEX
				accounts_account_name_active_account_uniq_idx ON accounts (account_name)
			WHERE account_status != 'deleted' AND account_status != 'banned';
			`,
		},
		Down: []string{
			`
			DROP INDEX accounts_account_name_not_null_uniq_idx;
			`,
		},
	},
}

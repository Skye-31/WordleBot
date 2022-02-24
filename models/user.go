package models

import "github.com/DisgoOrg/snowflake"

type User struct {
	ID                snowflake.Snowflake `bun:"id,pk,nullzero"`
	CachedTag         string              `bun:"tag,nullzero,notnull"`
	Public            bool                `bun:"public,notnull,default:false"`
	DefaultWordLength uint8               `bun:"default_word_length,notnull,nullzero,default:5"`
}

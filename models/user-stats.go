package models

import "github.com/disgoorg/snowflake"

type UserStats struct {
	ID    snowflake.Snowflake `bun:"id,pk,nullzero"`
	Four  []int               `bun:"four,type:smallint[]"`
	Five  []int               `bun:"five,type:smallint[]"`
	Six   []int               `bun:"six,type:smallint[]"`
	Seven []int               `bun:"seven,type:smallint[]"`
	Eight []int               `bun:"eight,type:smallint[]"`
}

package handlers

import "github.com/mateoschiro8/morfeo/server/types"

var TC TokenControler

type TokenControler struct {
	Tokens *map[string]*types.UserInput
}

func (tc *TokenControler) GetToken(id string) *types.UserInput {
	return (*tc.Tokens)[id]
}

func LoadTokenControler(tokens *map[string]*types.UserInput){
	
	if TC == (TokenControler{}){
		TC = TokenControler{
			Tokens : tokens,
		}
	}
}
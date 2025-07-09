package model

import "github.com/m1crogravity/spy-cat-agency/internal/validator"

const AgentUserType = UserType("agent")

var AnonymousAgent = &Agent{}

type Agent struct {
	Id       int64    `json:"id"`
	Name     string   `json:"name"`
	Password password `json:"-"`
}

func (a *Agent) IsAnonymous() bool {
	return a == AnonymousAgent
}

func ValidateAgent(v *validator.Validator, agent *Agent) {
	v.Check(agent.Name != "", "name", "must be provided")
	v.Check(len(agent.Name) <= 500, "name", "must be more than 500 bytes long")

	if agent.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *agent.Password.plaintext)
	}

	if agent.Password.hash == nil {
		panic("missing password hash for user")
	}
}

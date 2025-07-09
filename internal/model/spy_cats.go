package model

import "github.com/m1crogravity/spy-cat-agency/internal/validator"

var AnonymousSpyCat = &SpyCat{}

const SpyCatUserType = UserType("spy-cat")

type SpyCat struct {
	Id                int64    `json:"id"`
	Name              string   `json:"name"`
	YearsOfExperience uint     `json:"years_of_experience"`
	Breed             string   `json:"breed"`
	Salary            float64  `json:"salary"`
	Password          password `json:"-"`
}

func (sc *SpyCat) IsAnonymous() bool {
	return sc == AnonymousSpyCat
}

func ValidateSpyCat(v *validator.Validator, spyCat *SpyCat) {
	v.Check(spyCat.Name != "", "name", "must be provided")
	v.Check(len(spyCat.Name) <= 500, "name", "must be more than 500 bytes long")

	if spyCat.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *spyCat.Password.plaintext)
	}

	if spyCat.Password.hash == nil {
		panic("missing password hash for user")
	}
}

package model

type SpyCat struct {
	Id                int64   `json:"id"`
	Name              string  `json:"name"`
	YearsOfExperience uint    `json:"years_of_experience"`
	Breed             string  `json:"breed"`
	Salary            float64 `json:"salary"`
}

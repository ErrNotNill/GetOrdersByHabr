package models

type Resume struct {
	Position         string `json:"position,omitempty"`
	Company          string `json:"company,omitempty"`
	Salary           string `json:"salary,omitempty"`
	Schedule         string `json:"schedule,omitempty"`
	DescriptionTasks string `json:"description_tasks,omitempty"`
	LinkToVacancy    string `json:"ink_to_vacancy,omitempty"`
}

type Vacancy struct {
	Position          string `json:"position,omitempty"`
	Salary            string `json:"salary,omitempty"`
	Schedule          string `json:"schedule,omitempty"`
	DescriptionSkills string `json:"description_skills,omitempty"`
	FullName          string `json:"full_name,omitempty"`
	LinkToVacancy     string `json:"link_to_vacancy,omitempty"`
}

type Company struct {
	Inn  string `json:"inn"`
	Kpp  string `json:"kpp"`
	Name string `json:"name"`
	Boss string `json:"boss"`
}

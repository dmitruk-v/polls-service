package scheme

type Poll struct {
	SurveyID     int64             `json:"survey_id"`
	PreSetValues map[string]string `json:"pre_set_values"`
}

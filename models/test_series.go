package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type TestSeries struct {
	Id          int64         `orm:"auto"`
	Course      *Course       `orm:"rel(fk); null"`
	Name        string        `orm:"size(300); null" valid:"MaxSize(300)"`
	Description string        `orm:"type(text); null"`
	File        string        `orm:"size(400); null"`
	Questions   string        `orm:"type(jsonb); notnull"`
	Timer       time.Duration `orm:"notnull"`
	CreatedAt   time.Time     `orm:"auto_now_add;type(datetime)"`
	UpdatedAt   time.Time     `orm:"auto_now;type(datetime)"`
}

type TestSeriesSolution struct {
	Id          int64       `orm:"auto"`
	TestSeries  *TestSeries `orm:"rel(fk); notnull"`
	Description string      `orm:"type(text); null"`
	Solution    string      `orm:"type(jsonb); notnull"`
	CreatedAt   time.Time   `orm:"auto_now_add;type(datetime)"`
	UpdatedAt   time.Time   `orm:"auto_now;type(datetime)"`
}

type TestSeriesAttempt struct {
	Id         int64       `orm:"auto"`
	TestSeries *TestSeries `orm:"rel(fk); notnull" valid:"Required;"`
	User       *User       `orm:"rel(fk); notnull" valid:"Required;"`
	Result     string      `orm:"type(jsonb); notnull" valid:"Required;"`
	Submitted  bool        `orm:"default(false)"`
	TotalScore float64     `orm:"decimal(2)"`
	TotalMarks float64     `orm:"decimal(2)"`
	CreatedAt  time.Time   `orm:"auto_now_add;type(datetime)"`
	UpdatedAt  time.Time   `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(TestSeries))
	orm.RegisterModel(new(TestSeriesSolution))
	orm.RegisterModel(new(TestSeriesAttempt))
}

type ResultSerialzer struct {
	Solutions  *TestSeriesSolution `json:"solutions"`
	UserAnswer *TestSeriesAttempt  `json:"answer"`
}

func CreateResult(payload *SubmitTestSerializer, user *User) (*ResultSerialzer, error) {
	attempt := TestSeriesAttempt{}
	o := orm.NewOrm()
	err := o.QueryTable("test_series_attempt").
		Filter("TestSeries__id", payload.TestSeries).
		Filter("User__id", user).
		One(&attempt)

	if err != nil {
		return nil, errors.New("User attempt not found")
	}

	solutions := TestSeriesSolution{}

	errs := o.QueryTable("test_series_solution").
		Filter("TestSeries__id", payload.TestSeries).
		One(&solutions)

	if errs != nil {
		return nil, errors.New("TestSeries Solution not found")
	}

	return &ResultSerialzer{
		Solutions:  &solutions,
		UserAnswer: &attempt,
	}, nil
}

func CreateAttempt(payload *TestSeriesAttempt) (*TestSeriesAttempt, error) {
	o := orm.NewOrm()
	attempt := TestSeriesAttempt{}

	testSeriesExists := o.QueryTable("test_series").Filter("Id", payload.TestSeries.Id).Exist()

	if !testSeriesExists {
		return nil, errors.New("TestSeries is Not Available!")
	}

	// Try to fetch an existing attempt
	err := o.QueryTable("test_series_attempt").
		Filter("TestSeries__id", payload.TestSeries).
		Filter("User__id", payload.User).
		One(&attempt)

	var newResult map[string]interface{}
	if err := json.Unmarshal([]byte(payload.Result), &newResult); err != nil {
		return nil, errors.New("invalid result format, expected JSON object")
	}

	if err == orm.ErrNoRows {

		resultArray := []map[string]interface{}{newResult}
		resultJSON, _ := json.Marshal(resultArray)

		newAttempt := TestSeriesAttempt{
			TestSeries: payload.TestSeries,
			User:       payload.User,
			Result:     string(resultJSON), // Convert array to JSON string
		}

		// Insert into DB
		id, insertErr := o.Insert(&newAttempt)
		if insertErr != nil {
			return nil, insertErr
		}

		newAttempt.Id = id
		return &newAttempt, nil
	} else if err != nil {

		return nil, err
	}

	var resultArray []map[string]interface{}

	if attempt.Result != "" {
		if err := json.Unmarshal([]byte(attempt.Result), &resultArray); err != nil {
			return nil, errors.New("failed to parse existing result data")
		}
	}

	questionExists := false
	for i, item := range resultArray {
		if item["Question"] == newResult["Question"] {
			resultArray[i]["Answer"] = newResult["Answer"] // Update answer
			questionExists = true
			break
		}
	}

	if !questionExists {
		resultArray = append(resultArray, newResult) // Add new question-answer pair
	}

	updatedResultJSON, _ := json.Marshal(resultArray)

	attempt.Result = string(updatedResultJSON)
	_, updateErr := o.Update(&attempt, "Result", "UpdatedAt")
	if updateErr != nil {
		return nil, updateErr
	}

	return &attempt, nil
}

type SubmitTestSerializer struct {
	TestSeries int64 `json:"testseries" valid:"Required"`
}

func SubmitTestSeries(payload *SubmitTestSerializer, user *User) (*TestSeriesAttempt, error) {
	o := orm.NewOrm()
	attempt := TestSeriesAttempt{}

	testSeries := TestSeries{}
	err := o.QueryTable("test_series").Filter("Id", payload.TestSeries).One(&testSeries)
	if err != nil {
		return nil, err
	}

	solution := TestSeriesSolution{}
	err = o.QueryTable("test_series_solution").Filter("TestSeries__id", payload.TestSeries).One(&solution)
	if err != nil {
		return nil, errors.New("Solution not found for this test series")
	}

	var correctAnswers []map[string]interface{}
	if err := json.Unmarshal([]byte(solution.Solution), &correctAnswers); err != nil {
		return nil, errors.New("Failed to parse correct answers")
	}

	err = o.QueryTable("test_series_attempt").
		Filter("TestSeries__id", payload.TestSeries).
		Filter("User__id", user).
		One(&attempt)

	if err != nil {
		return nil, errors.New("User attempt not found")
	}

	var userAnswers []map[string]interface{}
	if err := json.Unmarshal([]byte(attempt.Result), &userAnswers); err != nil {
		return nil, errors.New("Failed to parse user answers")
	}

	userAnswersMap := make(map[string]string)
	for _, answer := range userAnswers {
		questionID := fmt.Sprintf("%v", answer["Question"])
		userAnswersMap[questionID] = fmt.Sprintf("%v", answer["Answer"])
	}

	totalScore := 0.0
	totalMarks := 0.0

	for _, question := range correctAnswers {
		questionID := fmt.Sprintf("%v", question["UUID"])
		answer := question["Answer"]
		positiveMark := question["Positive Marks"].(string)
		negativeMark := question["Negative Marks"].(string)
		positiveMarks, err := strconv.ParseFloat(positiveMark, 64)
		if err != nil {
			fmt.Println("Error converting Positive Marks:", err)
			positiveMarks = 0.0 // Default value if conversion fails
		}
		totalMarks += positiveMarks

		negativeMarks, err := strconv.ParseFloat(negativeMark, 64)
		if err != nil {
			fmt.Println("Error converting Positive Marks:", err)
			negativeMarks = 0.0 // Default value if conversion fails
		}

		if userAnswers, found := userAnswersMap[questionID]; found {
			if userAnswers == answer {
				totalScore += positiveMarks
			} else {
				totalScore -= negativeMarks
			}
		}
	}
	attempt.Submitted = true
	attempt.TotalMarks = totalMarks
	attempt.TotalScore = totalScore

	_, updateErr := o.Update(&attempt, "TotalMarks", "TotalScore", "Submitted", "UpdatedAt")
	if updateErr != nil {
		return nil, updateErr
	}

	return &attempt, nil
}

func CreateAttempts(payload *TestSeriesAttempt) (*TestSeriesAttempt, error) {
	o := orm.NewOrm()
	attempt := TestSeriesAttempt{}

	// ✅ Ensure TestSeries exists
	testSeriesExists := o.QueryTable("test_series").Filter("Id", payload.TestSeries).Exist()
	if !testSeriesExists {
		return nil, errors.New("TestSeries is Not Available!")
	}

	// ✅ Check if attempt exists
	err := o.QueryTable("test_series_attempt").
		Filter("TestSeries__id", payload.TestSeries).
		Filter("User__id", payload.User).
		One(&attempt)

	if err == orm.ErrNoRows {
		// 🔥 First attempt, insert new JSONB data
		insertQuery := `
			INSERT INTO test_series_attempt (test_series_id, user_id, result)
			VALUES (?, ?, ?::jsonb) RETURNING id
		`
		var newId int64

		insertErr := o.Raw(insertQuery, payload.TestSeries, payload.User, payload.Result).QueryRow(&newId)
		if insertErr != nil {
			return nil, insertErr
		}

		attempt.Id = newId
		attempt.TestSeries = payload.TestSeries
		attempt.User = payload.User
		attempt.Result = payload.Result
		return &attempt, nil
	} else if err != nil {
		return nil, err
	}

	// ✅ Extract Question from payload.Result
	extractQuestionQuery := `SELECT result->>'Question' FROM jsonb_array_elements(?::jsonb)`
	var question string
	err = o.Raw(extractQuestionQuery, payload.Result).QueryRow(&question)
	if err != nil {
		return nil, errors.New("Failed to extract Question from JSONB")
	}

	// ✅ Update `Answer` if `Question` exists, else append new Question-Answer
	updateQuery := `
		UPDATE test_series_attempt
		SET result = (
			CASE 
				WHEN EXISTS (
					SELECT 1 FROM jsonb_array_elements(result) elem 
					WHERE elem->>'Question' = ?
				) THEN (
					SELECT jsonb_agg(
						CASE 
							WHEN elem->>'Question' = ? THEN jsonb_set(elem, '{Answer}', to_jsonb(?))
							ELSE elem 
						END
					) 
					FROM jsonb_array_elements(result) elem
				)
				ELSE result || ?::jsonb
			END
		) WHERE id = ?
		RETURNING result
	`

	var updatedResult string
	updateErr := o.Raw(updateQuery, question, question, payload.Result, payload.Result, attempt.Id).QueryRow(&updatedResult)
	if updateErr != nil {
		return nil, updateErr
	}

	attempt.Result = updatedResult
	return &attempt, nil
}

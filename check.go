package onfido

import (
	"bytes"
	"context"
	"encoding/json"
	"time"
)

// CheckType represents a check type (express, standard)
type CheckType string

// CheckStatus represents a status of a check
type CheckStatus string

// CheckResult represents a result of a check (clear, consider)
type CheckResult string

// Supported check types
const (
	CheckStatusInProgress        CheckStatus = "in_progress"
	CheckStatusAwaitingApplicant CheckStatus = "awaiting_applicant"
	CheckStatusComplete          CheckStatus = "complete"
	CheckStatusWithdrawn         CheckStatus = "withdrawn"
	CheckStatusPaused            CheckStatus = "paused"
	CheckStatusReopened          CheckStatus = "reopened"

	CheckResultClear    CheckResult = "clear"
	CheckResultConsider CheckResult = "consider"
)

// CheckRequest represents a check request to Onfido API
type CheckRequest struct {
	ApplicantID             string   `json:"applicant_id"`
	ReportNames             []string `json:"report_names"`
	RedirectURI             string   `json:"redirect_uri,omitempty"`
	Tags                    []string `json:"tags,omitempty"`
	SupressFormEmails       bool     `json:"suppress_form_emails,omitempty"`
	Async                   bool     `json:"asynchronous,omitempty"`
	ChargeApplicantForCheck bool     `json:"charge_applicant_for_check,omitempty"`
	// Consider is used for Sandbox Testing of multiple report scenarios.
	// see https://documentation.onfido.com/#sandbox-responses
	Consider              []ReportName `json:"consider,omitempty"`
	ApplicantProvidesData bool         `json:"applicant_provides_data"`
}

// Check represents a check in Onfido API
type Check struct {
	ID                    string      `json:"id,omitempty"`
	CreatedAt             *time.Time  `json:"created_at,omitempty"`
	Href                  string      `json:"href,omitempty"`
	Type                  CheckType   `json:"type,omitempty"`
	Status                CheckStatus `json:"status,omitempty"`
	Result                CheckResult `json:"result,omitempty"`
	DownloadURI           string      `json:"download_uri,omitempty"`
	FormURI               string      `json:"form_uri,omitempty"`
	RedirectURI           string      `json:"redirect_uri,omitempty"`
	ResultsURI            string      `json:"results_uri,omitempty"`
	Reports               []*Report   `json:"reports,omitempty"`
	Tags                  []string    `json:"tags,omitempty"`
	ApplicantID           string      `json:"applicant_id,omitempty"`
	ApplicantProvidesData bool        `json:"applicant_provides_data"`
}

// CheckRetrieved represents a check in the Onfido API which has been retrieved.
// This is subtly different to the Check type above, as the Reports slice
// is just a string of Report IDs, not fully expanded Report objects.
// See https://documentation.onfido.com/?shell#check-object (Shell)
type CheckRetrieved struct {
	ID                    string      `json:"id,omitempty"`
	CreatedAt             *time.Time  `json:"created_at,omitempty"`
	Href                  string      `json:"href,omitempty"`
	Type                  CheckType   `json:"type,omitempty"`
	Status                CheckStatus `json:"status,omitempty"`
	Result                CheckResult `json:"result,omitempty"`
	DownloadURI           string      `json:"download_uri,omitempty"`
	FormURI               string      `json:"form_uri,omitempty"`
	RedirectURI           string      `json:"redirect_uri,omitempty"`
	ResultsURI            string      `json:"results_uri,omitempty"`
	Reports               []string    `json:"report_ids,omitempty"`
	Tags                  []string    `json:"tags,omitempty"`
	ApplicantID           string      `json:"applicant_id,omitempty"`
	ApplicantProvidesData bool        `json:"applicant_provides_data"`
}

// Checks represents a list of checks in Onfido API
type Checks struct {
	Checks []*Check `json:"checks"`
}

// CreateCheck creates a new check for the provided applicant.
// see https://documentation.onfido.com/?shell#create-check
func (c *client) CreateCheck(ctx context.Context, cr CheckRequest) (*Check, error) {
	jsonStr, err := json.Marshal(cr)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest("POST", "/checks", bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	var resp Check
	_, err = c.do(ctx, req, &resp)
	return &resp, err
}

// GetCheck retrieves a check by its ID
// see https://documentation.onfido.com/?shell#retrieve-check
func (c *client) GetCheck(ctx context.Context, id string) (*CheckRetrieved, error) {
	req, err := c.newRequest("GET", "/checks/"+id, nil)
	if err != nil {
		return nil, err
	}

	var resp CheckRetrieved
	_, err = c.do(ctx, req, &resp)
	return &resp, err
}

// GetCheckExpanded retrieves a check by its ID, with
// the Check's Reports expanded within the returned Check object.
// see https://documentation.onfido.com/?shell#retrieve-check (Shell) but refer to the JSON
// response object for https://documentation.onfido.com/?php#check-object (PHP) for the expanded contents.
func (c *client) GetCheckExpanded(ctx context.Context, id string) (*Check, error) {
	// Get the CheckRetrieved object. This only includes Report IDs, not the expanded Report objects.
	chkRetrieved, err := c.GetCheck(ctx, id)
	if err != nil {
		return nil, err
	}

	// Build a regular Check object, this is what will be returned assuming there is no error.
	check := Check{
		ApplicantID: chkRetrieved.ApplicantID,
		CreatedAt:   chkRetrieved.CreatedAt,
		DownloadURI: chkRetrieved.DownloadURI,
		FormURI:     chkRetrieved.FormURI,
		Href:        chkRetrieved.Href,
		ID:          chkRetrieved.ID,
		RedirectURI: chkRetrieved.RedirectURI,
		Reports:     make([]*Report, len(chkRetrieved.Reports)),
		Result:      chkRetrieved.Result,
		ResultsURI:  chkRetrieved.ResultsURI,
		Status:      chkRetrieved.Status,
		Tags:        chkRetrieved.Tags,
		Type:        chkRetrieved.Type,
	}

	// For each Report ID in the CheckRetrieved object, fetch (expand) the Report
	// into the returned Check object.
	for i, reportID := range chkRetrieved.Reports {
		rep, err := c.GetReport(ctx, reportID)
		if err != nil {
			return nil, err
		}
		check.Reports[i] = rep
	}
	return &check, nil
}

// ResumeCheck resumes a paused check by its ID.
// see https://documentation.onfido.com/?shell#resume-check
func (c *client) ResumeCheck(ctx context.Context, id string) (*Check, error) {
	req, err := c.newRequest("POST", "/checks/"+id+"/resume", nil)
	if err != nil {
		return nil, err
	}

	var resp Check
	_, err = c.do(ctx, req, &resp)
	return &resp, err
}

// CheckIter represents a check iterator
type CheckIter struct {
	*iter
}

// Check returns the current item in the iterator as a Check.
func (i *CheckIter) Check() *Check {
	return i.Current().(*Check)
}

// ListChecks retrieves the list of checks for the provided applicant.
// see https://documentation.onfido.com/?shell#list-checks
func (c *client) ListChecks(applicantID string) *CheckIter {
	handler := func(body []byte) ([]interface{}, error) {
		var r Checks
		if err := json.Unmarshal(body, &r); err != nil {
			return nil, err
		}

		values := make([]interface{}, len(r.Checks))
		for i, v := range r.Checks {
			values[i] = v
		}
		return values, nil
	}

	return &CheckIter{&iter{
		c:       c,
		nextURL: "/checks?applicant_id=" + applicantID,
		handler: handler,
	}}
}

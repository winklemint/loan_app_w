package lead

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

type LeadStore interface {
	InsertLead(lead LeadInfo) error
	GetLead(id int) (LeadInfo, error)
	UpdateLead(id int, lead LeadInfo) error
	DeleteLead(id int) error
	// GetLeadsByPincode(pincode int) ([]LeadInfo, error)

}

type MyLeadStore struct {
	Leads map[int]LeadInfo
}

func (s *MyLeadStore) InsertLead(lead LeadInfo) error {
	lead.ID = len(s.Leads) + 1
	lead.Created_at = time.Now()
	lead.Last_modified = lead.Created_at
	s.Leads[lead.ID] = lead
	return nil
}

func (s *MyLeadStore) GetLead(id int) (LeadInfo, error) {
	lead, exists := s.Leads[id]
	if !exists {
		return LeadInfo{}, errors.New("lead not found")
	}
	return lead, nil
}

func (s *MyLeadStore) UpdateLead(id int, lead LeadInfo) error {
	_, exists := s.Leads[id]
	if !exists {
		return errors.New("lead not found")
	}
	lead.ID = id
	lead.Last_modified = time.Now()
	s.Leads[id] = lead
	return nil
}

func (s *MyLeadStore) DeleteLead(id int) error {
	_, exists := s.Leads[id]
	if !exists {
		return errors.New("lead not found")
	}
	delete(s.Leads, id)
	return nil
}

// func (s *MyLeadStore) GetLeadsByPincode(pincode int) ([]LeadInfo, error) {
// 	var leads []LeadInfo
// 	for _, lead := range s.Leads {
// 		if lead.Pincode == pincode {
// 			leads = append(leads, lead)
// 		}
// 	}
// 	return leads, nil
// }

type LoanValidator struct {
	store LeadStore
}

func NewLoanValidator(store LeadStore) *LoanValidator {
	return &LoanValidator{
		store: store,
	}
}

func (v *LoanValidator) ValidateLoan(leadID int) error {
	lead, err := v.store.GetLead(leadID)
	if err != nil {
		return err
	}
	if lead.Loan_type == "" {
		return errors.New("Loan type is required.")
	}

	if lead.Loan_amount <= 0 {
		return errors.New("Loan amount must be greater than zero.")
	}

	if lead.Tenure <= 0 {
		return errors.New("Tenure must be greater than zero.")
	}

	if match, _ := regexp.MatchString(`^\d{6}$`, fmt.Sprint(lead.Pincode)); !match {
		return errors.New("Pincode must be exactly 6 digits.")
	}

	if lead.Employment_type == "" {
		return errors.New("Employment type is required.")
	}

	if lead.Gross_monthly_income <= 5000 {
		return errors.New("Gross monthly income must be greater than zero.")
	}

	return nil
}

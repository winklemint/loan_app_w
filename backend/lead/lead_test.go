package lead

import (
	"testing"
	//"time"
	//"time"
)

func TestInsertLead(t *testing.T) {
	store := &MyLeadStore{
		Leads: make(map[int]LeadInfo),
	}

	lead := LeadInfo{
		Loan_type:            "Home Loan",
		Loan_amount:          200000,
		Tenure:               24,
		Pincode:              123456,
		Employment_type:      "Full-time",
		Gross_monthly_income: 5000,
	}
	err := store.InsertLead(lead)
	if err != nil {
		t.Errorf("InsertLead returned an error: %v", err)
	}

	// Verify that the lead was inserted
	if len(store.Leads) != 1 {
		t.Errorf("Expected 1 lead, got %d", len(store.Leads))
	}
}

func TestGetLead(t *testing.T) {
	store := &MyLeadStore{
		Leads: make(map[int]LeadInfo),
	}

	lead := LeadInfo{
		Loan_type:            "Home Loan",
		Loan_amount:          200000,
		Tenure:               24,
		Pincode:              123456,
		Employment_type:      "Full-time",
		Gross_monthly_income: 5000,
	}
	err := store.InsertLead(lead)
	if err != nil {
		t.Errorf("InsertLead returned an error: %v", err)
	}

	// Retrieve the lead by ID
	retrievedLead, err := store.GetLead(1)
	if err != nil {
		t.Errorf("GetLead returned an error: %v", err)
	}

	// Verify the retrieved lead matches the inserted lead
	if retrievedLead.Loan_type != lead.Loan_type ||
		retrievedLead.Loan_amount != lead.Loan_amount ||
		retrievedLead.Tenure != lead.Tenure ||
		retrievedLead.Pincode != lead.Pincode ||
		retrievedLead.Employment_type != lead.Employment_type ||
		retrievedLead.Gross_monthly_income != lead.Gross_monthly_income {
		t.Errorf("Retrieved lead does not match the inserted lead")
	}
}

func TestUpdateLead(t *testing.T) {
	store := &MyLeadStore{
		Leads: make(map[int]LeadInfo),
	}

	lead := LeadInfo{
		Loan_type:            "Home Loan",
		Loan_amount:          200000,
		Tenure:               24,
		Pincode:              123456,
		Employment_type:      "Full-time",
		Gross_monthly_income: 5000,
	}
	err := store.InsertLead(lead)
	if err != nil {
		t.Errorf("InsertLead returned an error: %v", err)
	}

	updatedLead := LeadInfo{
		ID:                   1,
		Loan_type:            "Car Loan",
		Loan_amount:          150000,
		Tenure:               36,
		Pincode:              654321,
		Employment_type:      "Part-time",
		Gross_monthly_income: 4000,
	}
	err = store.UpdateLead(1, updatedLead)
	if err != nil {
		t.Errorf("UpdateLead returned an error: %v", err)
	}

	// Retrieve the updated lead
	retrievedLead, err := store.GetLead(1)
	if err != nil {
		t.Errorf("GetLead returned an error: %v", err)
	}

	// Verify the retrieved lead matches the updated lead
	if retrievedLead.Loan_type != updatedLead.Loan_type ||
		retrievedLead.Loan_amount != updatedLead.Loan_amount ||
		retrievedLead.Tenure != updatedLead.Tenure ||
		retrievedLead.Pincode != updatedLead.Pincode ||
		retrievedLead.Employment_type != updatedLead.Employment_type ||
		retrievedLead.Gross_monthly_income != updatedLead.Gross_monthly_income {
		t.Errorf("Retrieved lead does not match the updated lead")
	}
}

func TestDeleteLead(t *testing.T) {
	store := &MyLeadStore{
		Leads: make(map[int]LeadInfo),
	}

	lead := LeadInfo{
		Loan_type:            "Home Loan",
		Loan_amount:          200000,
		Tenure:               24,
		Pincode:              123456,
		Employment_type:      "Full-time",
		Gross_monthly_income: 5000,
	}
	err := store.InsertLead(lead)
	if err != nil {
		t.Errorf("InsertLead returned an error: %v", err)
	}

	err = store.DeleteLead(1)
	if err != nil {
		t.Errorf("DeleteLead returned an error: %v", err)
	}

	if len(store.Leads) != 0 {
		t.Errorf("Expected 0 leads, got %d", len(store.Leads))
	}
}

func TestLoanValidator_ValidateLoan(t *testing.T) {
	mockStore := &MyLeadStore{
		Leads: make(map[int]LeadInfo),
	}

	// Insert test leads
	mockStore.InsertLead(LeadInfo{
		Loan_type:            "Home Loan",
		Loan_amount:          200000,
		Tenure:               24,
		Pincode:              123456,
		Employment_type:      "Full-time",
		Gross_monthly_income: 25000,
	})

	mockStore.InsertLead(LeadInfo{
		Loan_type:            "Home Loan",
		Loan_amount:          200000,
		Tenure:               24,
		Pincode:              123456,
		Employment_type:      "Full-time",
		Gross_monthly_income: 25000,
	})

	validator := NewLoanValidator(mockStore)

	// Test case: Valid loan
	err := validator.ValidateLoan(1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test case: Invalid loan (income below threshold)
	err = validator.ValidateLoan(2)
	if err == nil {
		t.Error("Expected error, but got nil")
	} else {
		expectedErr := "income does not meet minimum threshold"
		if err.Error() != expectedErr {
			t.Errorf("Expected error message: %s, but got: %v", expectedErr, err)
		}
	}

	// Test case: Non-existent lead
	// err = validator.ValidateLoan(3)
	// if err == nil {
	// 	t.Error("Expected error, but got nil")
	// } else {
	// 	expectedErr := "lead not found"
	// 	if err.Error() != expectedErr {
	// 		t.Errorf("Expected error message: %s, but got: %v", expectedErr, err)
	// 	}
	// }
}

package loan

import (
	"fmt"
	"testing"
)

func TestInsertLoan(t *testing.T) {
	mockStore := &MockLoanStore{
		Loans:    make(map[int]Loan_details),
		NextID:   0,
		IsClosed: false,
	}

	//  using the mock_loan_store.go file for the insert function
	err := mockStore.InsertLoan(Loan_details{ID: 1, Loan_type: "home", Loan_amount: 60000, Tenure: 4, Pincode: 123456, Employment_type: "Salaried", Gross_monthly_income: 5000})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestUpdateLoan(t *testing.T) {
	mockStore := &MockLoanStore{
		Loans:    make(map[int]Loan_details),
		NextID:   0,
		IsClosed: false,
	}

	//inserting 1 detail of loan
	err := mockStore.InsertLoan(Loan_details{ID: 1, Loan_type: "home", Loan_amount: 60000, Tenure: 4, Pincode: 123456, Employment_type: "Salaried", Gross_monthly_income: 5000})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Updating the loan here after inserting
	err = mockStore.UpdateLoan(1, Loan_details{ID: 1, Loan_type: "car", Loan_amount: 50000, Tenure: 3, Pincode: 654321, Employment_type: "Self-employed", Gross_monthly_income: 6000})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// fetching ht updated loan details
	updatedLoan, err := mockStore.GetLoan(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fmt.Println(updatedLoan)

}

func TestDeleteLoan(t *testing.T) {
	mockStore := &MockLoanStore{
		Loans:    make(map[int]Loan_details),
		NextID:   0,
		IsClosed: false,
	}

	// for every function we are first inserting then following up with the function required in this case insertion then deletion
	err := mockStore.InsertLoan(Loan_details{ID: 1, Loan_type: "home", Loan_amount: 60000, Tenure: 4, Pincode: 123456, Employment_type: "Salaried", Gross_monthly_income: 5000})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Delete query
	err = mockStore.DeleteLoan(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// fetching the deleted loan if the delete function at the top gets failed then it will return the error
	_, err = mockStore.GetLoan(1)
	if err == nil {
		t.Fatalf("expected error: loan should be deleted")
	}

}

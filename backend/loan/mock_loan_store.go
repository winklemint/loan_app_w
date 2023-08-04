package loan

import (
	"errors"
)

type MockLoanStore struct {
	Loans    map[int]Loan_details
	loans    map[string]Loan_details
	NextID   int
	IsClosed bool
}

func (m *MockLoanStore) InsertLoan(loan Loan_details) error {

	if m.IsClosed {
		return errors.New("mock store is closed")
	}

	m.NextID++
	loan.ID = m.NextID
	m.Loans[loan.ID] = loan
	return nil
}

func (m *MockLoanStore) GetLoan(id int) (Loan_details, error) { // we are retriving the loan based on ID here
	if m.IsClosed {
		return Loan_details{}, errors.New("mock store is closed")
	}

	loan, exists := m.Loans[id]
	if !exists {
		return Loan_details{}, errors.New("loan not found")
	}

	return loan, nil
}

func (m *MockLoanStore) UpdateLoan(id int, loan Loan_details) error {
	if m.IsClosed {
		return errors.New("mock store is closed")
	}

	if _, exists := m.Loans[id]; !exists {
		return errors.New("loan not found")
	}

	loan.ID = id //update based on id
	m.Loans[id] = loan
	return nil
}

func (m *MockLoanStore) DeleteLoan(id int) error {
	if m.IsClosed {
		return errors.New("mock store is closed")
	}

	if _, exists := m.Loans[id]; !exists { // deleteion only if id exists
		return errors.New("loan not found")
	}

	delete(m.Loans, id)
	return nil
}

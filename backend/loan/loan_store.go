package loan

type LoanStore interface {
	InsertLoan(loan Loan_details) error
	GetLoan(id string) (Loan_details, error)
	UpdateLoan(id string, loan Loan_details) error
	DeleteLoan(id string) error
}

import React ,{useContext} from 'react'
import { MultiStepContext } from '../StepContext';

 const StepSeven = () => {
    const {setCurrentStep , userData, SetUserData,setSubmit} = useContext(MultiStepContext);

    const handleNext = () => {
      setCurrentStep(7);


    let loanDetails = { "loan_type" : userData.loan_type , "loan_amount" : parseFloat(userData.loanAmount), "pincode" : parseInt(userData.pincode), "employment_type" : userData.employment_type, "gross_monthly_income" : parseFloat(userData.GrossMonthly), "tenure" : parseFloat(userData.tenure) };
    console.log(loanDetails);

    fetch('proxy?url=http://127.0.0.1:8080/lead/add', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        loan_type: userData.loan_type,
        loan_amount: parseFloat(userData.loanAmount),
        pincode: parseInt(userData.pincode),
        employment_type: userData.employment_type,
        gross_monthly_income: parseFloat(userData.GrossMonthly),
        tenure : parseInt(userData.tenure)
      }),
    })
      .then((response) => {
        console.log(response);

        if (response.ok) {
          console.log(response.data);
          return response;
        } else {
          return response.json().then((data) => {
            let errorMessage = 'Authentication Failed';
            if (data && data.error && data.error.message) {
              errorMessage = data.error.message;
            }
            throw new Error(errorMessage);
          });
        }
      })
      .then((data) => {
        console.log('this is important', data);
      })
      .catch((error) => {
        console.log('There was an error', error);
      });
  

    }
    
  return (
    <div>
       <section>
        <div className="container d-flex justify-content-center align-items-center mt-5 box">
          <div className="row">
            <div className="card d-flex shadow-lg " style={{ backgroundColor: '#F7F8FA' }}>
              <div className="card-body">
                <div className="container text-center my-3">
                  <h2>
                    <b>Loan Application</b>
                  </h2>
                </div>

                <div className="firstbox w-100">
            <h2 className="fs-4">Step 6: tenure <p>(How Much Amount of time you want Loan)</p></h2>
            <div className="row">
            </div>

            <div className="form-group col-md-12 mt-3 hoverEffect">
              <label> Select tenure Date</label>
              <input
                type="text"
                name="test"
                id="cb15"
                placeholder="Fill Date"
                className="form-control"
                value={userData['tenure']}
                onChange={(e) => SetUserData({...userData, 'tenure': e.target.value})}
              />
            </div>

            <div className="d-flex">
              <button
                className="btn btn-success m-3 mt-5 w-50"
                onClick={()=>setCurrentStep(5)}
              >
                Previous
              </button>
            
                <button className="btn btn-success m-3 mt-5 w-50"onClick={handleNext} > Submit </button>
            </div>
          </div>
          </div>
          </div>
          </div>
          </div>
          </section>
    </div>
  )
}
export default StepSeven

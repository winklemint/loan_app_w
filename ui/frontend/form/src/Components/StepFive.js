import React, { useContext} from 'react'
import { MultiStepContext } from '../StepContext';

const StepFive = () => {
    const {setCurrentStep , userData, SetUserData} = useContext(MultiStepContext);

    const handleNext = () => {
      setCurrentStep(6);

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
            <h2 className="fs-4">Step 5: Property Location and Pincode</h2>
            <div className="row">
            </div>

            <div className="form-group col-md-12 mt-3 hoverEffect">
              <label> Enter Pincode:</label>
              <input
                type="text"
                name="test"
                id="cb15"
                placeholder="Enter Pincode"
                className="form-control"
                value={userData['pincode']}
                onChange={(e) => SetUserData({...userData, 'pincode': e.target.value})}
              />
            </div>

            <div className="d-flex">
              <button
                className="btn btn-success m-3 mt-5 w-50"
                onClick={()=>setCurrentStep(4)}
              >
                Previous
              </button>
            
                <button className="btn btn-success m-3 mt-5 w-50" onClick={handleNext} > Next </button>
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
export default StepFive
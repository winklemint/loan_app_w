import React, {useState , useContext} from 'react'
import { MultiStepContext } from '../StepContext';

 const StepSix = () =>{
    const {setCurrentStep , userData, SetUserData,setSubmit} = useContext(MultiStepContext);
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
                  {/* <div className="progress mb-4 firstbox">
                    <div
                      className="prog progress-bar" role="progressbar"
                      style={{ width: `${calculateProgress()}%`, color: `` }}
                      aria-valuenow={calculateProgress()}
                      aria-valuemin="0"
                      aria-valuemax="100"
                    ></div>
                  </div> */}
                </div>

                <div className="firstbox w-100">
            <h1 className="subm">Successfully submitted!</h1>
            <div className='text-center'>
            <button className="btn btn-success m-3 mt-5"onClick={()=>setSubmit} >See Details </button>
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
export default StepSix

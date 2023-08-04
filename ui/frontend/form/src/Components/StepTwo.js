import React, { useState, useContext, useEffect } from 'react'
import { MultiStepContext } from '../StepContext';


const StepTwo = () => {
    const { setCurrentStep, userData, SetUserData, saveData } = useContext(MultiStepContext);

    const [selectedOption, setSelectedOption] = useState(userData.employment_type || '');
    const [isOptionSelected, setIsOptionSelected] = useState(false);

    const handleRadioChange = (event) => {
        setSelectedOption(event.target.value);
        setIsOptionSelected(true);
    };

    useEffect(() => {
        SetUserData({ ...userData, employment_type: selectedOption });
    },[selectedOption])

     

        const handleNext = () => {
            if (isOptionSelected) {
                const updatedData = {...userData, employment_type:selectedOption};
                SetUserData(updatedData);
              // Save the selected loan purpose to the database
                saveData(updatedData)
                .then(() => {
                    setCurrentStep(3);
                })
              // Proceed to the next step
            
              .catch((error) => {
                console.error(
                  'There has been a problem with your fetch operation:',
                  error
                );
                alert('An error occurred while saving data. Please try again.');
              });
            } else {
              alert('Please select a employment type.');
            }
          };
      
    return (
        <div>
            <section>
                <div className="container d-flex justify-content-center align-items-center mt-5 box">
                    <div className="row">
                        <div
                            className="card d-flex shadow-lg "
                            style={{ backgroundColor: '#F7F8FA' }}
                        >
                            <div className="card-body">
                                <div className="container my-3">
                                    <h2>
                                        <b>Loan Application</b>
                                    </h2>
                                </div>
                                <div>
                                    <div className="firstbox w-100">
                                        <h2 className="fs-4">Step 2: Employement Type</h2>
                                        <div className="row">
                                            <div className="form-group col-md-6 hoverEffect">
                                                <input type="radio" name="test" id="cb8" className="form-control" value="salaried"
                                                    checked={selectedOption === 'salaried'}
                                                    onChange={handleRadioChange} />
                                                <label for="cb8" className="checked-label">
                                                    <img src="assets/money.png" alt='eight' />
                                                    Salaried
                                                </label>
                                            </div>

                                            <div className="form-group col-md-6 hoverEffect">
                                                <input
                                                    type="radio"
                                                    name="test"
                                                    id="cb9"
                                                    className="form-control"
                                                    value="BussinessOwner"
                                                    checked={selectedOption === 'BussinessOwner'}
                                                    onChange={handleRadioChange}
                                                />
                                                <label for="cb9" className="checked-label">
                                                    <img src="assets/mortgage-loan.png" alt='nine' />
                                                    Bussiness Owner
                                                </label>
                                            </div>

                                            <div className="form-group col-md-6 hoverEffect">
                                                <input
                                                    type="radio"
                                                    name="test"
                                                    id="cb10"
                                                    className="form-control"
                                                    value="SelfEmployed"
                                                    checked={selectedOption === 'SelfEmployed'}
                                                    onChange={handleRadioChange}
                                                />
                                                <label for="cb10" className="checked-label">
                                                    <img src="assets/self-employed.png" alt='ten' />
                                                    Self Employeed
                                                </label>
                                            </div>

                                            <div className="form-group col-md-6 hoverEffect">
                                                <input
                                                    type="radio"
                                                    name="test"
                                                    id="cb11"
                                                    className="form-control"
                                                    value="Independent"
                                                    checked={selectedOption === 'Independent'}
                                                    onChange={handleRadioChange}
                                                />
                                                <label for="cb11" className="checked-label">
                                                    <img src="assets/worker.png" alt='eleven' />
                                                    Independent Worker
                                                </label>
                                            </div>
                                        </div>

                                        <div className="text-center mt-5" style={{ color: 'gray' }}>
                                            <p>Select Employement type</p>
                                        </div>

                                        <div className="d-flex">
                                            <button className="btn btn-success m-3  w-50" onClick={() => setCurrentStep(1)}>
                                                Previous
                                            </button>

                                            <button
                                                className="btn btn-success m-3  w-50"
                                               onClick={handleNext}
                                            >
                                                Next
                                            </button>
                                        </div>
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
export default StepTwo

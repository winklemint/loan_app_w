import React, { useState,useContext }from 'react'
import { MultiStepContext } from '../StepContext';



const StepTHree = () => {

    const {setCurrentStep , userData, SetUserData, saveData} = useContext(MultiStepContext);

    const [value, setValue] = useState(userData.loanAmount || ''); // Initial value of the slider

              const [Range, setRange] = useState('');
            
              const convertToWords = (num) => {
                const units = [
                  '',
                  'one',
                  'two',
                  'three',
                  'four',
                  'five',
                  'six',
                  'seven',
                  'eight',
                  'nine',
                ];
                const tens = [
                  '',
                  '',
                  'twenty',
                  'thirty',
                  'forty',
                  'fifty',
                  'sixty',
                  'seventy',
                  'eighty',
                  'ninety',
                ];
                const thousands = ['', 'thousand', 'lakh'];
            
                if (num === 0) return 'zero';
            
                const convertChunkToWords = (chunk) => {
                  let chunkWords = '';
            
                  const hundreds = Math.floor(chunk / 100);
                  const remainder = chunk % 100;
            
                  if (hundreds > 0) {
                    chunkWords += units[hundreds] + ' hundred';
                  }
            
                  if (remainder > 0) {
                    if (chunkWords !== '') {
                      chunkWords += ' ';
                    }
            
                    if (remainder < 20) {
                      chunkWords += units[remainder];
                    } else {
                      const tensPlace = Math.floor(remainder / 10);
                      const onesPlace = remainder % 10;
            
                      chunkWords += tens[tensPlace];
            
                      if (onesPlace > 0) {
                        chunkWords += '-' + units[onesPlace];
                      }
                    }
                  }
            
                  return chunkWords;
                };
            
                let words = '';
                let chunkIndex = 0;
            
                while (num > 0) {
                  const chunk = num % 1000;
                  if (chunk !== 0) {
                    const chunkWords = convertChunkToWords(chunk);
                    if (chunkIndex > 0) {
                      words = chunkWords + ' ' + thousands[chunkIndex] + ' ' + words;
                    } else {
                      words = chunkWords;
                    }
                  }
                  num = Math.floor(num / 1000);
                  chunkIndex++;
                }
                return words;
              };

              const handleChange = (event) => {
                setValue(event.target.value);
                SetUserData((prevUserData) => ({ ...prevUserData, loanAmount: event.target.value }));
              };

              const handleNext = async () => {
                if (value) {
                  // Proceed to the next step
                  SetUserData((prevUserData) => ({ ...prevUserData, loanAmount: value }));

                 
            try {
              // Save the updated data to the server
              await saveData({ ...userData, loanAmount: value });

              // Proceed to the next step
              setCurrentStep(4);
          } catch (error) {
              console.error('There has been a problem with your fetch operation:', error);
              alert('An error occurred while saving data. Please try again.');
          }
            
                } else {
                  alert('Please select a loan amount.');
                }
              };

              const words = convertToWords(value);
              
            
            
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

            <div className="firstbox w-100">
            <h2 className="fs-4">Step 3: Loan Amount Required</h2>
            <div className="row">
              <input
                // type="range"
                min={0}
                max={1000000}
                value={value}
               onChange={handleChange}
              />
               <p>Value: {value}</p>
              <p>Words: {words}</p>
            </div>
            <div className="d-flex">
              <button
                className="btn btn-success m-3 mt-5 w-50"
                onClick={() => setCurrentStep(2)}
              >
                Previous
              </button>
              
                <button
                  className="btn btn-success m-3 mt-5 w-50"
                  onClick={handleNext} >  Next</button>
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
export default  StepTHree

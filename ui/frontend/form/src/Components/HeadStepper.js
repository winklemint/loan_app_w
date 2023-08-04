import React, { useContext } from 'react';
import '../App.css';
import Navbar from './Navbar';
import StepOne from './StepOne';
import StepTwo from './StepTwo';
import StepTHree from './StepTHree';
import StepFour from './StepFour';
import StepFive from './StepFive';
import StepSix from './StepSix';
import StepSeven from './StepSeven';
import { Stepper, StepLabel, Step } from '@mui/material';
import { MultiStepContext } from '../StepContext';

function HeadStepper() {
  const { currentStep } = useContext(MultiStepContext);
  const showStep = (step) => {
    console.log(step);
    switch (step) {
      case 1:
        return <StepOne />;

      case 2:
        return <StepTwo />;

      case 3:
        return <StepTHree />;

      case 4:
        return <StepFour />;

      case 5:
        return <StepFive />;

      case 6:
          return <StepSeven />;

      case 7:
        return <StepSix />;
    }
  };
  return (
    <>
      <Navbar />

      <div className="center-stepper mt-3">
        <Stepper
          style={{ width: '60%', position: 'relative', left: '20%' }}
          activeStep={currentStep - 1}
          orientation="horizontal"
        >
          <Step>
            <StepLabel></StepLabel>
          </Step>

          <Step>
            <StepLabel></StepLabel>
          </Step>

          <Step>
            <StepLabel></StepLabel>
          </Step>

          <Step>
            <StepLabel></StepLabel>
          </Step>

          <Step>
            <StepLabel></StepLabel>
          </Step>

          <Step>
            <StepLabel></StepLabel>
          </Step>

          <Step>
            <StepLabel></StepLabel>
          </Step>
        </Stepper>
      </div>
      {showStep(currentStep)}
    </>
  );
}

export default HeadStepper;
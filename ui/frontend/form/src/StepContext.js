import React, { useState, createContext } from 'react';
import App from './App';

export const saveData = async (userData) => {
  try{
  const response = await fetch('http://127.0.0.1:8080/loan/insert', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(userData),
  });

  if (!response.ok) {
    const message = await response.json();
    throw new Error(message);
  }
  const data = await response.json();
  return data;
} catch (error) {
  console.error("Error in saveData:", error);
  throw error;
}
};

export const MultiStepContext = createContext();

const StepContext = () => {
  const [currentStep, setCurrentStep] = useState(1);
  const [userData, SetUserData] = useState([]);
  const [finalData, SetFinalData] = useState([]);

  return (
    <div>
      <MultiStepContext.Provider value={{ currentStep, setCurrentStep, userData, SetUserData, finalData, SetFinalData, saveData }}>
        <App/>
      </MultiStepContext.Provider>
    </div>
  );
}

export default StepContext;

import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import SignUp from './Components/signup';
import Login from './Components/login';
import './App.css';
import '../node_modules/bootstrap/dist/css/bootstrap.min.css';
import ForgotPass from './Components/forgetpass';
import HeadStepper from './Components/HeadStepper';
import Otp from './Components/Otp';
import NewPassword from './Components/NewPassword';
import Mail from './Components/Mail';
import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import ProtectedRoute from './Components/ProtectedRoute';

function App() {
  return (
    <div className="App">
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<SignUp />} />
          <Route path="/Login" element={<Login />} />
          <Route path="/Forgotpass" element={<ForgotPass />} />
          <Route path="/otp" element={<Otp />} />
          <Route path="/NewPassword" element={<NewPassword />} />
          <Route path="/Mail" element={<Mail />} />
          {/* Protected routes should be wrapped inside a single <Route> */}
          <Route
            path="/protected"
            element={
              <>
                <ProtectedRoute path="/HeadStepper" element={<HeadStepper />} />
                {/* <ProtectedRoute path="/otp" element={<Otp />} /> */}
              </>
            }
          />
        </Routes>
      </BrowserRouter>
      <ToastContainer />
    </div>
  );
}

export default App;

import React, { useState } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';

const OTPVerification = () => {
  const [otp, setOTP] = useState('');
  const [otpError, setOTPError] = useState('');
  const navigate = useNavigate();
  const location = useLocation();
  const queryParams = new URLSearchParams(location.search);
  const email = queryParams.get('email');

  const handleChange = (index, value) => {
    const newOTP = otp.split('');
    newOTP[index] = value;
    setOTP(newOTP.join(''));
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    setOTPError('');

    // Validate the OTP field
    if (!otp) {
      setOTPError('Please enter the OTP');
      return;
    }
    fetch('/proxy?url=http://127.0.0.1:8080/otpverify', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: email, // Replace with the user's email
        otp: otp,
      }),
    })
      .then((data) => {
        console.log('Response from server:', data);
        if (data.success) {
          setOTPError(data.error || 'Invalid OTP. Please try again.');
        } else {
          console.log('OTP Verification Successful');
          navigate('/NewPassword');
        }
      })
      .catch((error) => {
        console.error('Error during OTP verification', error);
        setOTPError(
          'An error occurred during OTP verification. Please try again later.'
        );
      });
  };

  return (
    <div className="container h-100vh con1">
      <div className="row">
        <div className="col-md-2"></div>
        <div className="col-md-8">
          <div className="card d-flex auth-inner ">
            <div className="card py-5 px-3">
              <h5 className="m-0">OTP verification</h5>
              <span className="mobile-text">
                Enter the code we just sent to your Email{' '}
                <b className="text-danger">{email}</b>
              </span>
              <div className="d-flex flex-row mt-5">
                {/* Use the otp state value and handleChange function */}
                {[0, 1, 2, 3, 4, 5].map((index) => (
                  <input
                    key={index}
                    type="text"
                    className="form-control m-1"
                    value={otp.charAt(index)}
                    onChange={(e) => handleChange(index, e.target.value)}
                    maxLength={1}
                    autoFocus={index === 0} // Auto-focus the first input
                  />
                ))}
              </div>
              {/* Display OTP validation error message */}
              {otpError && <div className="error">{otpError}</div>}
              <div className="text-center mt-5">
                <span className="d-block mobile-text">
                  Didn't receive the code?
                </span>
                <h6 className="m-2">
                  <Link to="/Forgetpass">Login</Link>
                </h6>
              </div>
              {/* Form submit button */}
              <div className="text-center mt-5">
                <button
                  type="submit"
                  className="btn btn-primary"
                  onClick={handleSubmit}
                >
                  Verify OTP
                </button>
              </div>
            </div>
          </div>
        </div>
        <div className="col-md-2"></div>
      </div>
    </div>
  );
};

export default OTPVerification;

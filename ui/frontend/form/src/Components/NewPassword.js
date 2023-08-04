// CreateNewPassword.js
import React, { useState } from 'react';
import { Link } from 'react-router-dom';

function CreateNewPassword() {
  // State variables for password fields and errors
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [passwordError, setPasswordError] = useState('');

  // Function to handle form submission
  const handleSubmit = (e) => {
    e.preventDefault();

    // Reset the error message
    setPasswordError('');

    // Validate the password fields
    if (!password) {
      setPasswordError('Please enter a password');
      return;
    }

    if (password !== confirmPassword) {
      setPasswordError('Passwords do not match');
      return;
    }
  };

  return (
    <div className="container h-100vh con1">
      <div className="row">
        <div className="col-md-2"></div>
        <div className="col-md-8">
          <div className="card d-flex auth-inner">
            <div className="card-body">
              <form className="" onSubmit={handleSubmit}>
                <h3>Create New Password</h3>
                <div className="mb-3 ">
                  <label>Password : </label>
                  <input
                    type="password"
                    className="form-control"
                    placeholder="Enter password"
                    name="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                  />
                </div>
                <div className="mb-3 ">
                  <label>Confirm Password : </label>
                  <input
                    type="password"
                    className="form-control"
                    placeholder="Confirm password"
                    name="confirmPassword"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    required
                  />
                </div>
                {passwordError && <div className="error">{passwordError}</div>}
                <button type="submit" className="btn btn-primary ">
                  Update Password
                </button>{' '}
                <br />
                <Link to="/">Back to Login</Link>
              </form>
            </div>
          </div>
        </div>
        <div className="col-md-2"></div>
      </div>
    </div>
  );
}

export default CreateNewPassword;

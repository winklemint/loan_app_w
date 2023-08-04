import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';

function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [contactNum, setContactNum] = useState('');
  const [errors, setErrors] = useState({});
  const navigate = useNavigate();

  const handleSubmit = (e) => {
    e.preventDefault();

    let formErrors = {};
    if (!email) {
      formErrors.email = 'Please enter your email';
    } else if (!validateEmail(email)) {
      formErrors.email = 'Please enter a valid email';
    }
    if (!contactNum) {
      formErrors.contactNum = 'Please enter your phone contact number';
    } else if (!validatePhone(contactNum)) {
      formErrors.contactNum = 'Please enter a valid phone number';
    }
    if (!password) {
      formErrors.password = 'Please enter your password';
    } else if (password.length < 8) {
      formErrors.password = 'Password must contain at least 8 characters';
    }
    setErrors(formErrors);

    if (
      Object.keys(formErrors).length === 0 &&
      email &&
      contactNum &&
      password
    ) {
      fetch('/proxy?url=http://127.0.0.1:8080/user/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: email,
          contact_num: parseInt(contactNum),
          password: password,
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
              //   console.log(errorData);
              throw new Error(errorMessage);
            });
          }
        })
        .then((data) => {
          console.log('this is important', data);
          navigate('/HeadStepper');
        })
        .catch((error) => {
          console.log('There was an error', error);
        });
    }
  };

  const validateEmail = (email) => {
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
  };

  const validatePhone = (contactNum) => {
    return /^[0-9]{10}$/.test(contactNum);
  };

  return (
    <div className="container h-100vh con1">
      <div className="row">
        <div className="col-md-2"></div>
        <div className="col-md-8">
          <div className="card d-flex auth-inner">
            <div className="card-body">
              <form className="needs-validation" onSubmit={handleSubmit}>
                <h3>Login</h3>

                <div className="mb-3 ">
                  <label>Email</label>
                  <input
                    type="email"
                    className="form-control"
                    placeholder="Enter email"
                    name="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                  />
                  {errors.email && (
                    <span className="error" style={{ color: 'red' }}>
                      {errors.email}
                    </span>
                  )}
                </div>

                <div className="mb-3 ">
                  <label>Phone</label>
                  <input
                    type="tel"
                    name="contactNum"
                    value={contactNum}
                    className="form-control"
                    placeholder="Enter phone number"
                    onChange={(e) => setContactNum(e.target.value)}
                  />
                  {errors.contactNum && (
                    <span className="error" style={{ color: 'red' }}>
                      {errors.contactNum}
                    </span>
                  )}
                </div>

                <div className="mb-3 ">
                  <label>Password</label>
                  <input
                    type="password"
                    className="form-control"
                    placeholder="Enter password"
                    name="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                  />
                  {errors.password && (
                    <span className="error" style={{ color: 'red' }}>
                      {errors.password}
                    </span>
                  )}
                </div>
                <Link to="/Forgotpass">Forgot Password</Link>

                <div className="d-grid">
                  <button type="submit" className="btn btn-primary">
                    Login
                  </button>
                </div>
                <p className="forgot-password text-right">
                  Don't have an account? <Link to="/">Sign in</Link>
                </p>
              </form>
            </div>
          </div>
        </div>
        <div className="col-md-2"></div>
      </div>
    </div>
  );
}

export default Login;

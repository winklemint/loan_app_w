import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';

function SignUp() {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [contact, setcontact] = useState('');
  const [password, setPassword] = useState('');
  const [errors, setErrors] = useState({});
  const navigate = useNavigate();

  const handleSubmit = (e) => {
    e.preventDefault();

    let formErrors = {};
    if (!name) {
      formErrors.name = 'Please enter your name';
    }
    if (!email) {
      formErrors.email = 'Please enter your email';
    } else if (!validateEmail(email)) {
      formErrors.email = 'Please enter a valid email';
    }
    if (!password) {
      formErrors.password = 'Please enter your password';
    } else if (password.length < 8) {
      formErrors.password = 'Password must contain at least 8 characters';
    }
    if (!contact) {
      formErrors.contact = 'Please enter your phone contact';
    } else if (!validatePhone(contact)) {
      formErrors.contact = 'Please enter a valid phone number';
    }

    setErrors(formErrors);

    if (
      Object.keys(formErrors).length === 0 &&
      name &&
      email &&
      contact &&
      password
    ) {
      fetch('/proxy?url=http://127.0.0.1:8080/user/add', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name: name,
          email: email,
          contact: parseInt(contact),
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
              throw new Error(errorMessage);
            });
          }
        })
        .then((data) => {
          console.log('Registration Successfull', data);
          navigate('/login');
        })
        .catch((error) => {
          console.error('Error during registration', error);
        });
    }
  };

  const validateEmail = (email) => {
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
  };

  const validatePhone = (contact) => {
    return /^[0-9]{10}$/.test(contact);
  };

  return (
    <div className="container h-100vh con1">
      <div className="row">
        <div className="col-md-2"></div>
        <div className="col-md-8">
          <div className="card d-flex auth-inner">
            <div className="card-body">
              <form className="needs-validation" onSubmit={handleSubmit}>
                <h3>Registration Form</h3>
                <div className="mb-3">
                  <label>Name</label>
                  <input
                    type="text"
                    value={name}
                    className="form-control"
                    placeholder="Enter Name"
                    name="name"
                    onChange={(e) => setName(e.target.value)}
                  />
                  {errors.name && (
                    <span className="error" style={{ color: 'red' }}>
                      {errors.name}
                    </span>
                  )}
                </div>

                <div className="mb-3 ">
                  <label>Email</label>
                  <input
                    type="email"
                    value={email}
                    className="form-control"
                    placeholder="Enter email"
                    name="email"
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
                    name="contact"
                    value={contact}
                    className="form-control"
                    placeholder="Enter Phone contact"
                    onChange={(e) => setcontact(e.target.value)}
                  />
                  {errors.contact && (
                    <span className="error" style={{ color: 'red' }}>
                      {errors.contact}
                    </span>
                  )}
                </div>

                <div className="mb-3 ">
                  <label>Password</label>
                  <input
                    type="password"
                    className="form-control"
                    placeholder="Enter password"
                    value={password}
                    name="password"
                    onChange={(e) => setPassword(e.target.value)}
                  />
                  {errors.password && (
                    <span className="error" style={{ color: 'red' }}>
                      {errors.password}
                    </span>
                  )}
                </div>

                <div className="d-grid">
                  <button type="submit" className="btn btn-primary">
                    Sign Up
                  </button>
                </div>
                <p className="forgot-password text-right">
                  Already registered <Link to="/login">Log in?</Link>
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

export default SignUp;

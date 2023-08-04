import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import 'bootstrap/dist/css/bootstrap.min.css';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

function ForgotPass() {
  // State variables for form fields and errors
  const [email, setEmail] = useState('');
  const [emailError, setEmailError] = useState('');
  const navigate = useNavigate();

  // Function to handle form submission
  const handleSubmit = (e) => {
    e.preventDefault();

    let formErrors = {};

    if (!email) {
      formErrors.email = 'Please enter your Email';
    }

    setEmailError(formErrors);

    if (Object.keys(formErrors).length === 0 && email) {
      fetch('/proxy?url=http://127.0.0.1:8080/forgotpassword', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: email,
        }),
      })
        .then((response) => {
          console.log(response);

          if (response.ok) {
            toast.success('OTP sent successfully!', {
              autoClose: 5000, // Duration in milliseconds
              position: toast.POSITION.TOP_CENTER,
              hideProgressBar: false,
              closeOnClick: true,
              pauseOnHover: true,
              draggable: true,
            });
            console.log(response.data);
            return response;
          } else {
            return response.json().then((data) => {
              let errorMessage = 'Email not matched';
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
          navigate('/otp');
        })
        .catch((error) => {
          console.log('There was an error', error);
          toast.error('Error: Enter valid email', {
            position: toast.POSITION.TOP_CENTER,
            marginBottom: '70px',
            autoClose: 5000, // Duration in milliseconds
            hideProgressBar: false,
            closeOnClick: true,
            pauseOnHover: true,
            draggable: true,
          });
        });
    }
  };

  return (
    <div className="container h-100vh con1">
      <div className="row">
        <div className="col-md-2"></div>
        <div className="col-md-8">
          <div className="card d-flex auth-inner">
            <div className="card-body">
              <form className="needs-validation" onSubmit={handleSubmit}>
                <h3>Forgot Password</h3>
                {/* <img
                  src="https://cdn4.iconfinder.com/data/icons/emojis-flat-pixel-perfect/64/emoji-56-512.png"
                  alt="react logo"
                  className=""
                  style={{
                    width: '100px',
                    marginLeft: '75px',
                    border: ' 8px black',
                    animation: 'mymove 5s infinite',
                  }}
                /> */}
                <p>We will send you One time password on your Email</p>
                <div className="mb-3 was-validated">
                  <label>Email : </label>
                  <input
                    type="email"
                    className="form-control"
                    placeholder="Enter email"
                    name="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                  />
                  {emailError.email && (
                    <span className="error" style={{ color: 'red' }}>
                      {emailError.email}
                    </span>
                  )}
                </div>

                <button
                  type="submit"
                  className="btn btn-primary"
                  position="right center"
                >
                  Reset Password
                </button>
                <br />

                {/* Use the handleVerifyOTP function to navigate */}
                {/* <button className="btn btn-link">Verify OTP</button> */}
                <h6 className="m-2">
                  <Link to="/Login">Login</Link>
                </h6>
              </form>
            </div>
          </div>
        </div>
        <div className="col-md-2"></div>
      </div>
      <ToastContainer />
    </div>
  );
}

export default ForgotPass;

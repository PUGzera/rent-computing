import React, {useState} from 'react';
import { Formik, Form, Field, ErrorMessage } from 'formik';
import { Link, useNavigate } from 'react-router-dom';
import Header from '../header/Header';
 
function Register() {
    const [error, setError] = useState('')
    const navigate = useNavigate();
    return (
    <div>
        <Header />
        <div className="background" aria-hidden="true">
            <div className="cube"></div>
            <div className="cube"></div>
            <div className="cube"></div>
            <div className="cube"></div>
            <div className="cube"></div>
            <div className="cube"></div>
        </div>
        <div className='form-container'>
            <h1>Register</h1>
            <Formik
            initialValues={{ email: '', username: '', password: '' }}
            validate={values => {
                const errors = {};
                if (!values.email) {
                    errors.email = 'Required';
                }
                if (!values.username) {
                    errors.username = 'Required';
                }
                if (!values.password) {
                    errors.password = 'Required';
                }
                if (!values.confirmPassword) {
                    errors.confirmPassword = 'Required';
                }
                if (values.password !== values.confirmPassword) {
                    errors.confirmPassword = 'Passwords do not match';
                }
                return errors;
            }}
            onSubmit={(values, { setSubmitting }) => {
                setSubmitting(true)
                fetch("http://localhost:8080/users/register", {
                    method: "POST",
                    body: JSON.stringify({
                        "email": values.email,
                        "username": values.username,
                        "password": values.password
                    }),
                    headers: {
                        "Content-type": "application/json; charset=UTF-8"
                    }
                }).then((res) => {
                    if (!res.ok) {
                        throw new Error(`HTTP error! status: ${res.status}`);
                    }
                    return res.json();
                })
                .then((json) => {
                    console.log(json)
                    navigate("/login")
                })
                .catch((err) => setError(err.message))
            }}
            >
            {({ isSubmitting }) => (
                <Form>
                    <div className='form-group'>
                        <label htmlFor='email'>Email</label>
                        <Field name="email" type="email" />
                        <ErrorMessage name="email" component="div" />
                    </div>
                    <div className='form-group'>
                        <label htmlFor='username'>Username</label>
                        <Field name="username" type="text" />
                        <ErrorMessage name="username" component="div" />
                    </div>
                    <div className='form-group'>
                        <label htmlFor='password'>Password</label>
                        <Field type="password" name="password" />
                        <ErrorMessage name="password" component="div" />
                    </div>
                    <div className='form-group'>
                        <label htmlFor='confirmPassword'>Confirm Password</label>
                        <Field type="password" name="confirmPassword" />
                        <ErrorMessage name="confirmPassword" component="div" />
                    </div>
                    <div className='form-group'>
                        <button type="submit" disabled={isSubmitting}>
                            Register
                        </button>
                    </div>
                </Form>
            )}
            </Formik>
            <div style={{
                backgroundColor: '#ffe5e5',
                color: '#d8000c',
                border: '1px solid #f5c6cb',
                padding: '15px',
                borderRadius: '6px',
                marginTop: '15px',
                display: 'flex',
                alignItems: 'center',
                boxShadow: '0 2px 5px rgba(0,0,0,0.1)',
                display: error ? 'contents' : 'none'
                }}>
                <span style={{ marginRight: '10px' }}>⚠️</span>
                <span><strong>Error:</strong> {error}</span>
            </div>
            <div>
                <h3>Already have an account, sign in!</h3>
                <Link to="/login" className="btn secondary">
                    Sign In
                </Link>
            </div>
        </div>
    </div>
    )
}
 
 export default Register;
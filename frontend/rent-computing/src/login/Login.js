import React, {useState} from 'react';
import { Formik, Form, Field, ErrorMessage } from 'formik';
import { Link, useNavigate } from 'react-router-dom';
import Header from '../header/Header';
 
function Login({ setAuthenticated }) {
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
            <h1>Sign In</h1>
            <Formik
            initialValues={{ username: '', password: '' }}
            validate={values => {
                const errors = {};
                if (!values.username) {
                    errors.username = 'Required';
                }
                if (!values.password) {
                    errors.password = 'Required';
                }
                return errors;
            }}
            onSubmit={(values, { setSubmitting }) => {
                setSubmitting(true)
                fetch("http://localhost:8080/users/login", {
                    method: "POST",
                    body: JSON.stringify({
                        "username": values.username,
                        "password": values.password
                    }),
                    headers: {
                        "Content-type": "application/json; charset=UTF-8"
                    },
                    credentials: "include"
                }).then((res) => {
                    if (!res.ok) {
                        throw new Error(`HTTP error! status: ${res.status}`);
                    }
                    return res.json();
                })
                .then((json) => {
                    console.log(json)
                    setAuthenticated(true)
                })
                .catch((err) => setError(err.message))
                setSubmitting(false)
            }}
            >
            {({ isSubmitting }) => (
                <Form>
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
                        <button type="submit" disabled={isSubmitting} onClick={() => setError('')}>
                            Sign In
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
                <h3>Don't have an account, register!</h3>
                <Link to="/register" className="btn secondary">
                    Register
                </Link>
            </div>
        </div>
    </div>
    )
}
 
 export default Login;
import React, {useState} from 'react';
import { Formik, Form, Field, ErrorMessage } from 'formik';
import { Link, useNavigate } from 'react-router-dom';
import Header from '../header/Header';
import FormError from '../error/FormError';
 
function Register({ register }) {
    const [error, setError] = useState('')
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
                setError('')
                register(values.email, values.username, values.password, (_) => {setSubmitting(false); setError('')}, (err) => {setSubmitting(false); setError(err)})
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
            <FormError message={error} />
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
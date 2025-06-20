import React, {useState} from 'react';
import { Formik, Form, Field, ErrorMessage } from 'formik';
import { Link } from 'react-router-dom';
import Header from '../header/Header';
import FormError from '../error/FormError';
 
function Login({ login }) {
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
                setError('')
                login(values.username, values.password, (_) => {
                    setSubmitting(false) 
                    setError('')
                }, (err) => {
                    setSubmitting(false)
                    setError(err)
                })
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
                        <button type="submit" disabled={isSubmitting}>
                            {isSubmitting ? "Loading" : "Sign In" }
                        </button>
                    </div>
                </Form>
            )}
            </Formik>
            <FormError message={error} />
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
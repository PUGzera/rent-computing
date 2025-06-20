import Header from "./header/Header"
import { Formik, Form, Field, ErrorMessage } from 'formik';
import React, {useState, useEffect} from 'react';
import './profile.css'
import FormError from "../error/FormError";
import InfoToast from "../info/Info";
import ErrorToast from "../error/Error";

function Profile({ signOut, user, createMachine, listMachines }) {
    const [error, setError] = useState('')
    const [createFormError, setCreateFormError] = useState('')
    const [info, setInfo] = useState('')
    const [machines, setMachines] = useState([])

    useEffect(() => listMachines((json) => {
        setInfo('Succesfully fetched machines')
        setMachines(json)
    }, (err) => {
        setError(err)
    }), []);

    return (
        <div>
            <Header onSignOut={signOut} username={user.username}/>
            <div className="background" aria-hidden="true">
                <div className="cube"></div>
                <div className="cube"></div>
                <div className="cube"></div>
                <div className="cube"></div>
                <div className="cube"></div>
                <div className="cube"></div>
            </div>
            <div id="machines">
                <div className="machine-box">    
                    <img src="https://img.icons8.com/ios-filled/100/computer.png" alt="Computer" />
                    <Formik
                        initialValues={{ username: '', password: '' }}
                        validate={values => {
                            var errors = {}
                            if (!values.password) {
                                errors.password = 'Required';
                            }
                            return errors       
                        }}
                        onSubmit={(values, { setSubmitting }) => {
                            setSubmitting(true)
                            createMachine(values.password, (_) => {
                                setSubmitting(false) 
                                setCreateFormError('')
                                setInfo('Machine successfully created')
                            }, (err) => {
                                setSubmitting(false)
                                setCreateFormError(err)
                            })
                        }}
                        >
                        {({ isSubmitting }) => (
                            <Form>
                                <div className='form-group'>
                                    <label htmlFor='password'>Password</label>
                                    <Field type="password" name="password" />
                                    <ErrorMessage name="password" component="div" />
                                </div>
                                <div className='form-group'>
                                    <button type="submit" disabled={isSubmitting}>
                                        {isSubmitting ? "Loading" : "Create" }
                                    </button>
                                </div>
                            </Form>
                        )}
                    </Formik>
                    <FormError message={createFormError} />
                </div>
                {machines.map((machine) => (
                    <div className="machine-box">    
                        <img src="https://img.icons8.com/ios-filled/100/computer.png" alt="Computer" />
                        <p>Id: {machine.id}</p>
                        <button onClick={() => {window.open(`http://localhost${machine.address}`, '_blank')}}>Connect</button>
                    </div>
                ))}
            </div>
            <ErrorToast message={error} onClose={() => setError('')} />
            <InfoToast message={info} onClose={() => setInfo('')} />
        </div>
    )
}

export default Profile
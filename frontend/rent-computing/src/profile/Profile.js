import Header from "./header/Header"
import { Link, useNavigate } from 'react-router-dom';
import React, {useState} from 'react';

function Profile({ setAuthenticated }) {
    const [username, setUsername] = useState('')
    const navigate = useNavigate();
    const signOut = () => {
        fetch("http://localhost:8080/users/logout", {
            method: "GET",
            credentials: "include"
        }).then((res) => {
            if (res.ok) {
                setAuthenticated(false)
            }
        })
    }

    const getUsername = () => {
        fetch("http://localhost:8080/users/profile", {
            method: "GET",
            credentials: "include"
        }).then((res) => {
            if (!res.ok) {
                setUsername("COULD NOT FETCH USERNAME")
                throw new Error(`HTTP error! status: ${res.status}`);
            }
            return res.json()
        }).then((json) => {
            console.log(json)
            setUsername(json.user.username)
        })
    }

    getUsername()
    return (
        <div>
            <Header onSignOut={signOut} username={username}/>
            <div className="background" aria-hidden="true">
                <div className="cube"></div>
                <div className="cube"></div>
                <div className="cube"></div>
                <div className="cube"></div>
                <div className="cube"></div>
                <div className="cube"></div>
            </div>
        </div>
    )
}

export default Profile
import logo from './logo.svg';
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import React, {useState, useEffect} from 'react';
import './App.css';
import Login from './login/Login';
import Register from './register/Register';
import Home from './home/Home';
import Profile from './profile/Profile';
import ErrorToast from './error/Error';

function App() {
  const [authenticated, setAuthenticated] = useState(false)
  const [user, setUser] = useState(null)
  const [error, setError] = useState('')

  const isAuthenticated = () => {
    fetch("http://localhost:8080/users/profile", {
        method: "GET",
        credentials: "include"
    }).then((res) => {
        setAuthenticated(res.ok)
        if (!res.ok) {
          throw new Error(`HTTP error! status: ${res.status}`);
        }
        refresh()
        return res.json();
    }).then((json) => setUser(json.user))
    .catch((err) => setError(err.message))
  }

  const refresh = () => {
    fetch("http://localhost:8080/users/refresh", {
        method: "GET",
        credentials: "include"
    }).then((res) => {
        setAuthenticated(res.ok)
        if (!res.ok) {
          throw new Error(`HTTP error! status: ${res.status}`);
        }
    })
    .catch((err) => setError(err.message))
  }

  const login = (username, password, onSuccess, onFailure) => {
    fetch("http://localhost:8080/users/login", {
        method: "POST",
        body: JSON.stringify({
            "username": username,
            "password": password
        }),
        headers: {
            "Content-type": "application/json; charset=UTF-8"
        },
        credentials: "include"
    }).then((res) => {
        if (!res.ok) {
            throw res
        }
        return res.json();
    })
    .then((json) => {
        console.log(json)
        onSuccess(json)
        setAuthenticated(true)
    })
    .catch((err) => {
        if (err.text) {
            err.text().then(message => {
                var msg = `Status code ${err.status}, ${err.statusText}: Invalid credentials`
                setError(msg)
                onFailure(msg)
            })
        } else {
            setError(err.message)
            onFailure(err.message)
        }
    })
  }

  const register = (email, username, password, onSuccess, onFailure) => {
    fetch("http://localhost:8080/users/register", {
        method: "POST",
        body: JSON.stringify({
            "email": email,
            "username": username,
            "password": password
        }),
        headers: {
            "Content-type": "application/json; charset=UTF-8"
        }
    }).then((res) => {
        if (!res.ok) {
            throw res
        }
        return res.json();
    })
    .then((json) => {
        console.log(json)
        onSuccess(json)
        window.location = "/login"
    })
    .catch((err) => {
        if (err.text) {
            err.text().then(message => {
                var msg = `Status code ${err.status}, ${err.statusText}: ${JSON.parse(message).error}`
                setError(msg)
                onFailure(msg)
            })
        } else {
            setError(err.message)
            onFailure(err.message)
        }
    })
  }

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

  const createMachine = (password, onSuccess, onFailure) => {
    fetch("http://localhost:8080/machines/create", {
        method: "POST",
        body: JSON.stringify({
            "vnc": true,
            "password": password
        }),
        headers: {
            "Content-type": "application/json; charset=UTF-8"
        },
        credentials: "include"
    }).then((res) => {
        if (!res.ok) {
            throw res
        }
        return res.json();
    })
    .then((json) => {
        console.log(json)
        onSuccess(json)
    })
    .catch((err) => {
        if (err.text) {
            err.text().then(message => {
                var msg = `Status code ${err.status}, ${err.statusText}: ${JSON.parse(message).error}`
                setError(msg)
                onFailure(msg)
            })
        } else {
            setError(err.message)
            onFailure(err.message)
        }
    })
  }

  const listMachine = (onSuccess, onFailure) => {
    fetch("http://localhost:8080/machines/", {
        method: "GET",
        headers: {
            "Content-type": "application/json; charset=UTF-8"
        },
        credentials: "include"
    }).then((res) => {
        if (!res.ok) {
            throw res
        }
        return res.json();
    })
    .then((json) => {
        console.log(json)
        onSuccess(json)
    })
    .catch((err) => {
        if (err.text) {
            err.text().then(message => {
                var msg = `Status code ${err.status}, ${err.statusText}: ${JSON.parse(message).error}`
                setError(msg)
                onFailure(msg)
            })
        } else {
            setError(err.message)
            onFailure(err.message)
        }
    })
  }

  useEffect(() => isAuthenticated(setAuthenticated, setUser), []);
  
  if (authenticated) {
    return (
      <div>
        <ErrorToast message={error} onClose={() => setError('')} />
        <BrowserRouter>
          <Routes>
            <Route path="/profile" element={<Profile user={user} signOut={signOut} createMachine={createMachine} listMachines={listMachine} />} />
            <Route path="*" element={<Navigate to='/profile' replace />} />
          </Routes>
        </BrowserRouter>
      </div>
    )
  }
  return (
    <div>
      <ErrorToast message={error} onClose={() => setError('')} />
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/login" element={<Login login={login} />} /> 
          <Route path="/register" element={<Register register={register} />} /> 
          <Route path="*" element={<Navigate to='/' replace />} />
        </Routes>
      </BrowserRouter>
    </div>
  );
}

export default App;

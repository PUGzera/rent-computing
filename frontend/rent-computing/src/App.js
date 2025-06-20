import logo from './logo.svg';
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import React, {useState, useEffect} from 'react';
import './App.css';
import Login from './login/Login';
import Register from './register/Register';
import Home from './home/Home';
import Profile from './profile/Profile';

function authenticate(setAuthenticated) {
  fetch("http://localhost:8080/users/profile", {
      method: "GET",
      credentials: "include"
  }).then((res) => {
      setAuthenticated(res.ok)
  })
}

function App() {
  const [authenticated, setAuthenticated] = useState(false)
  useEffect(() => {
    authenticate(setAuthenticated);
  }, []);
  if (authenticated) {
    return (
      <BrowserRouter>
        <Routes>
          <Route path="/profile" element={<Profile setAuthenticated={setAuthenticated} />} />
          <Route path="*" element={<Navigate to='/profile' replace />} />
        </Routes>
      </BrowserRouter>
    )
  }
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/login" element={<Login setAuthenticated={setAuthenticated} />} /> 
        <Route path="/register" element={<Register />} /> 
        <Route path="*" element={<Navigate to='/' replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;

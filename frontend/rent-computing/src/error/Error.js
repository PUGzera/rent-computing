import React, { useEffect, useState } from 'react';
import './error.css';

const ErrorToast = ({ message, duration = 3000, onClose }) => {
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    if (message) {
      setVisible(true);
      const timer = setTimeout(() => {
        setVisible(false);
        setTimeout(onClose, 300); // Wait for fade-out before removing
      }, duration);
      return () => clearTimeout(timer);
    }
  }, [message, duration, onClose]);

  if (!message) return null;

  return (
    <div className={`error-toast ${visible ? 'show' : 'hide'}`}>
      <span className="error-toast-message">{message}</span>
      <button className="close-btn" onClick={() => {
        setVisible(false);
        setTimeout(onClose, 300); // manual close
      }}>
        &times;
      </button>
    </div>
  );
};

export default ErrorToast;
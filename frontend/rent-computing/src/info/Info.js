import React, { useEffect, useState } from 'react';
import './info.css';

const InfoToast = ({ message, duration = 3000, onClose }) => {
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
    <div className={`info-toast ${visible ? 'show' : 'hide'}`}>
      <span className="info-toast-message">{message}</span>
      <button className="close-info-btn" onClick={() => {
        setVisible(false);
        setTimeout(onClose, 300); // manual close
      }}>
        &times;
      </button>
    </div>
  );
};

export default InfoToast;
"use client";
import React, { useState } from 'react';
import { CSSProperties } from 'react';

const containerStyle: CSSProperties  = {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'flex-start',
    minHeight: '100vh',
    padding: '40px',
  };

const formStyle: CSSProperties  = {
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'flex-start',
  width: '500px',
};

const inputStyle: CSSProperties  = {
  width: '100%',
  padding: '16px 24px',
  margin: '12px 0',
  boxSizing: 'border-box',
  border: '1px solid #ccc',
  borderRadius: '4px',
  fontSize: '18px',
};

const labelStyle: CSSProperties  = {
  fontWeight: 'bold',
  color: '#ff6600',
  fontSize: '20px',
  marginBottom: '8px',
};

const fieldContainerStyle: CSSProperties  = {
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'flex-start',
  marginBottom: '24px',
};

const errorMessageStyle: CSSProperties  = {
  color: '#ff6600',
  marginBottom: '20px',
  fontSize: '16px',
};
const [isHovered, setIsHovered] = useState(false);

const buttonStyle: CSSProperties  = {
  width: '100%',
  backgroundColor: '#ff6600',
  color: 'white',
  padding: '18px 24px',
  margin: '12px 0',
  border: 'none',
  borderRadius: '4px',
  cursor: 'pointer',
  fontSize: '18px',
  boxShadow: '0 4px 8px rgba(0, 0, 0, 0.3)',
  transition: 'transform 0.2s ease-in-out',
  transform: isHovered ? 'scale(1.05)' : 'none', // Apply transform only when hovered
};
const buttonHoverStyle: CSSProperties  = {
  transform: 'scale(1.05)',
};

const titleStyle: CSSProperties  = {
  fontWeight: 'bold',
  fontSize: '32px',
  marginBottom: '24px',
  color: '#ff6600',
};

interface UserCredentials {
  email: string;
  password: string;
}

const Page = () => {
    const [formData, setFormData] = useState<UserCredentials>({
      email: '',
      password: '',
    });
    const [errorMessage, setErrorMessage] = useState('');
  
    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      const { name, value } = e.target;
      setFormData((prevData) => ({ ...prevData, [name]: value }));
    };
  
    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
      e.preventDefault();
      try {
        const response = await fetch('/users/login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(formData),
        });
        const data = await response.json();
        console.log(data);
        if (response.ok) {
          // Login successful, navigate to /home
          window.location.href = "/home";
        } else {
          // Invalid credentials
          setErrorMessage('Invalid credentials');
          setFormData({ email: '', password: '' }); // Clear form fields
        }
      } catch (error) {
        console.error('Error:', error);
      }
    };
  
    return (
      <div style={containerStyle}>
        <h1 style={titleStyle}>Login Form</h1>
        {errorMessage && (
          <div style={errorMessageStyle}>{errorMessage}</div>
        )}
        <form onSubmit={handleSubmit} style={formStyle}>
          <div style={fieldContainerStyle}>
            <label style={labelStyle}>Email</label>
            <input
              type="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              placeholder="Email"
              style={inputStyle}
            />
          </div>
          <div style={fieldContainerStyle}>
            <label style={labelStyle}>Password</label>
            <input
              type="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              placeholder="Password"
              style={inputStyle}
            />
          </div>
          <button 
        type="submit" 
        style={buttonStyle}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
      >
        Login
      </button>
        </form>
      </div>
    );
  };
  
  export default Page;
"use client";
import React, { useState } from 'react';

const containerStyle = {
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'flex-start',
  minHeight: '100vh',
  padding: '40px',
};

const formStyle = {
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'flex-start',
  width: '500px',
};

const inputStyle = {
  width: '100%',
  padding: '16px 24px',
  margin: '12px 0',
  boxSizing: 'border-box',
  border: '1px solid #ccc',
  borderRadius: '4px',
  fontSize: '18px',
};

const labelStyle = {
  fontWeight: 'bold',
  color: '#ff6600',
  fontSize: '20px',
  marginBottom: '8px',
};

const fieldContainerStyle = {
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'flex-start',
  marginBottom: '24px',
};

const errorMessageStyle = {
  color: '#ff6600',
  marginBottom: '20px',
  fontSize: '16px',
};

const buttonStyle = {
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
};

const buttonHoverStyle = {
  transform: 'scale(1.05)',
};

const titleStyle = {
  fontWeight: 'bold',
  fontSize: '32px',
  marginBottom: '24px',
  color: '#ff6600',
};

interface UserDetails {
  name: string;
  email: string;
  password: string;
}

const Page = () => {
  const [formData, setFormData] = useState<UserDetails>({
    name: '',
    email: '',
    password: '',
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prevData) => ({ ...prevData, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      const response = await fetch('/api/signup', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
      });
      const data = await response.json();
      console.log(data);
      if (response.ok) {
        // Signup successful, navigate to /home
        window.location.href = "/home";
      }
    } catch (error) {
      console.error('Error:', error);
    }
  };

  return (
    <div style={containerStyle}>
      <h1 style={titleStyle}>Signup Form</h1>
      <form onSubmit={handleSubmit} style={formStyle}>
        <div style={fieldContainerStyle}>
          <label style={labelStyle}>Name</label>
          <input
            type="text"
            name="name"
            value={formData.name}
            onChange={handleChange}
            placeholder="Name"
            style={inputStyle}
          />
        </div>
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
        <button type="submit" style={{ ...buttonStyle, ':hover': buttonHoverStyle }}>
          Signup
        </button>
      </form>
    </div>
  );
};

export default Page;
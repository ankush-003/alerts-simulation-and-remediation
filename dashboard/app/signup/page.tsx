"use client";
import React, { useState } from 'react';
import { CSSProperties } from 'react';

const containerStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'flex-start',
  minHeight: '100vh',
  padding: '40px',
};

const formStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'flex-start',
  width: '500px',
};

const inputStyle: CSSProperties = {
  width: '100%',
  padding: '16px 24px',
  margin: '12px 0',
  boxSizing: 'border-box',
  border: '1px solid #ccc',
  borderRadius: '4px',
  fontSize: '18px',
};

const labelStyle: CSSProperties = {
  fontWeight: 'bold',
  color: '#ff6600',
  fontSize: '20px',
  marginBottom: '8px',
};

const fieldContainerStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'flex-start',
  marginBottom: '24px',
};

const errorMessageStyle: CSSProperties = {
  color: '#ff6600',
  marginBottom: '20px',
  fontSize: '16px',
};

const titleStyle: CSSProperties = {
  fontWeight: 'bold',
  fontSize: '32px',
  marginBottom: '24px',
  color: '#ff6600',
};

interface UserDetails {
  First_name: string;
  Last_name: string;
  Email: string;
  Password: string;
  Phone: string;
  User_type: string;
}

const Page = () => {
  const [formData, setFormData] = useState<UserDetails>({
    First_name: '',
    Last_name: '',
    Email: '',
    Password: '',
    Phone: '', 
    User_type: 'USER'
  });

  const [isHovered, setIsHovered] = useState(false); // Move state declaration inside the component

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prevData) => ({ ...prevData, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      console.log(formData)
      const response = await fetch('http://localhost:8000/users/signup', {
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
        window.location.href = "/login";
      }
    } catch (error) {
      console.error('Error:', error);
    }
  };

  const buttonStyle: CSSProperties = {
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

  return (
    <div style={containerStyle}>
      <h1 style={titleStyle}>Signup Form</h1>
      <form onSubmit={handleSubmit} style={formStyle}>
        <div style={fieldContainerStyle}>
          <label style={labelStyle}>First Name</label>
          <input
            type="text"
            name="First_name"
            value={formData.First_name}
            onChange={handleChange}
            placeholder="First Name"
            style={inputStyle}
          />
        </div>
        <div style={fieldContainerStyle}>
          <label style={labelStyle}>Last Name</label>
          <input
            type="text"
            name="Last_name"
            value={formData.Last_name}
            onChange={handleChange}
            placeholder="Last Name"
            style={inputStyle}
          />
        </div>
        <div style={fieldContainerStyle}>
          <label style={labelStyle}>Email</label>
          <input
            type="email"
            name="Email"
            value={formData.Email}
            onChange={handleChange}
            placeholder="Email"
            style={inputStyle}
          />
        </div>
        <div style={fieldContainerStyle}>
          <label style={labelStyle}>Password</label>
          <input
            type="password"
            name="Password"
            value={formData.Password}
            onChange={handleChange}
            placeholder="Password"
            style={inputStyle}
          />
        </div>
        <div style={fieldContainerStyle}>
          <label style={labelStyle}>Phone</label>
          <input
            type="text"
            name="Phone"
            value={formData.Phone}
            onChange={handleChange}
            placeholder="Phone"
            style={inputStyle}
          />
        </div>
        <button 
          type="submit" 
          style={buttonStyle}
          onMouseEnter={() => setIsHovered(true)}
          onMouseLeave={() => setIsHovered(false)}
        >
          Signup
        </button>
      </form>
    </div>
  );
};

export default Page;

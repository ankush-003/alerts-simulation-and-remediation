"use client";
import React, { useState } from 'react';

const containerStyle: React.CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'flex-start',
  minHeight: '100vh',
  padding: '40px',
};

const formStyle: React.CSSProperties = {
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'flex-start',
  width: '900px',
  marginTop: '80px',
};

const inputStyle: React.CSSProperties = {
  width: '100%',
  padding: '16px 24px',
  margin: '12px 0',
  boxSizing: 'border-box',
  border: '1px solid #ccc',
  borderRadius: '4px',
  fontSize: '18px',
};

const labelStyle: React.CSSProperties = {
  fontWeight: 'bold',
  fontSize: '20px',
  marginBottom: '8px',
  display: 'flex',
  alignItems: 'center',
};

const fieldContainerStyle: React.CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'flex-start',
  marginBottom: '24px',
  marginRight: '48px',
};

const errorMessageStyle: React.CSSProperties = {
  color: '#ff6600',
  marginBottom: '20px',
  fontSize: '16px',
};

const titleStyle: React.CSSProperties = {
  fontWeight: 'bold',
  fontSize: '36px',
  marginBottom: '24px',
  color: '#ff6600',
};

export default function AlertConfig() {
  const [categories, setCategories] = useState({
    memory: false,
    cpu: false,
    disk: false,
    power: false,
  });
  const [severities, setSeverities] = useState({
    warning: false,
    critical: false,
    error: false,
  });

  const handleCategoryChange = (event) => {
    setCategories({ ...categories, [event.target.name]: event.target.checked });
  };

  const handleSeverityChange = (event) => {
    setSeverities({ ...severities, [event.target.name]: event.target.checked });
  };

  const handleSubmit = async () => {
    const alertConfig = {
      categories,
      severities,
    };

    const jwtToken = sessionStorage.getItem('token'); // Assuming the JWT token is stored in localStorage
    console.log(jwtToken, alertConfig)

    try {
      const response = await fetch('http://localhost:8000/users/alertconfig', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${jwtToken}`,
        },
        body: JSON.stringify(alertConfig),
      });

      if (response.ok) {
        console.log('Alert configuration saved successfully');
      } else {
        console.error('Failed to save alert configuration');
      }
    } catch (error) {
      console.error('Error saving alert configuration:', error);
    }
  };

  return (
    <div style={containerStyle}>
      <h2 style={titleStyle}>Alert Configuration</h2>
      <div style={formStyle}>
        <div style={fieldContainerStyle}>
          <h3 style={{ ...titleStyle, fontSize: '28px', color: '#ffff00' }}>Categories</h3>
          <div style={{ ...labelStyle, color: '#ffff00' }}>
            <input
              type="checkbox"
              name="memory"
              checked={categories.memory}
              onChange={handleCategoryChange}
              style={inputStyle}
            />
            <span style={{ marginLeft: '8px' }}>Memory</span>
          </div>
          <div style={{ ...labelStyle, color: '#ffff00' }}>
            <input
              type="checkbox"
              name="cpu"
              checked={categories.cpu}
              onChange={handleCategoryChange}
              style={inputStyle}
            />
            <span style={{ marginLeft: '8px' }}>CPU</span>
          </div>
          <div style={{ ...labelStyle, color: '#ffff00' }}>
            <input
              type="checkbox"
              name="disk"
              checked={categories.disk}
              onChange={handleCategoryChange}
              style={inputStyle}
            />
            <span style={{ marginLeft: '8px' }}>Disk</span>
          </div>
          <div style={{ ...labelStyle, color: '#ffff00' }}>
            <input
              type="checkbox"
              name="power"
              checked={categories.power}
              onChange={handleCategoryChange}
              style={inputStyle}
            />
            <span style={{ marginLeft: '8px' }}>Power</span>
          </div>
        </div>
        <div style={fieldContainerStyle}>
          <h3 style={{ ...titleStyle, fontSize: '28px', color: '#ff0000' }}>Severities</h3>
          <div style={{ ...labelStyle, color: '#ff0000' }}>
            <input
              type="checkbox"
              name="warning"
              checked={severities.warning}
              onChange={handleSeverityChange}
              style={inputStyle}
            />
            <span style={{ marginLeft: '8px' }}>Warning</span>
          </div>
          <div style={{ ...labelStyle, color: '#ff0000' }}>
            <input
              type="checkbox"
              name="critical"
              checked={severities.critical}
              onChange={handleSeverityChange}
              style={inputStyle}
            />
            <span style={{ marginLeft: '8px' }}>Critical</span>
          </div>
          <div style={{ ...labelStyle, color: '#ff0000' }}>
            <input
              type="checkbox"
              name="error"
              checked={severities.error}
              onChange={handleSeverityChange}
              style={inputStyle}
            />
            <span style={{ marginLeft: '8px' }}>Error</span>
          </div>
        </div>
      </div>
      <button onClick={handleSubmit} style={{ marginTop: '24px' }}>
        Submit
      </button>
      
    </div>
  );
}
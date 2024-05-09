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

const submitButtonStyle: React.CSSProperties = {
  backgroundColor: '#ff6600',
  color: 'white',
  padding: '12px 24px',
  border: 'none',
  borderRadius: '4px',
  fontSize: '18px',
  cursor: 'pointer',
  marginTop: '24px',
};


export default function AlertConfig() {
  let [categories, setCategories] = useState<any>([]);
  let [severities, setSeverities] = useState<any>([]);

  const handleCategoryChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const category = event.target.name;
    if (event.target.checked) {
      setCategories([...categories, category]);
    } else {
      setCategories(categories.filter((c:any) => c !== category));
    }
  };

  const handleSeverityChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const severity = event.target.name;
    if (event.target.checked) {
      setSeverities([...severities, severity]);
    } else {
      setSeverities(severities.filter((s:any) => s !== severity));
    }
  };

  const handleSubmit = async () => {
    let alertConfig = {
      categories,
      severities,
    };

    const jwtToken = sessionStorage.getItem('token');

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
        window.location.href = "/home"
      } else {
        const errorData = await response.json();
        console.error('Failed to save alert configuration', errorData.error);
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
              name="Memory"
              checked={categories.memory}
              onChange={handleCategoryChange}
              style={inputStyle}
            />
            <span style={{ marginLeft: '8px' }}>Memory</span>
          </div>
          <div style={{ ...labelStyle, color: '#ffff00' }}>
            <input
              type="checkbox"
              name="CPU"
              checked={categories.cpu}
              onChange={handleCategoryChange}
              style={inputStyle}
            />
            <span style={{ marginLeft: '8px' }}>CPU</span>
          </div>
          <div style={{ ...labelStyle, color: '#ffff00' }}>
            <input
              type="checkbox"
              name="Disk"
              checked={categories.disk}
              onChange={handleCategoryChange}
              style={inputStyle}
            />
            <span style={{ marginLeft: '8px' }}>Disk</span>
          </div>
          <div style={{ ...labelStyle, color: '#ffff00' }}>
            <input
              type="checkbox"
              name="Power"
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
              name="Warning"
              checked={severities.warning}
              onChange={handleSeverityChange}
              style={inputStyle}
            />
            <span style={{ marginLeft: '8px' }}>Warning</span>
          </div>
          <div style={{ ...labelStyle, color: '#ff0000' }}>
            <input
              type="checkbox"
              name="Critical"
              checked={severities.critical}
              onChange={handleSeverityChange}
              style={inputStyle}
            />
            <span style={{ marginLeft: '8px' }}>Critical</span>
          </div>
          <div style={{ ...labelStyle, color: '#ff0000' }}>
            <input
              type="checkbox"
              name="Severe"
              checked={severities.severe}
              onChange={handleSeverityChange}
              style={inputStyle}
            />
            <span style={{ marginLeft: '8px' }}>Severe</span>
          </div>
        </div>
      </div>
      <button onClick={handleSubmit} style={submitButtonStyle}>
        Submit
      </button>
      
    </div>
  );
}
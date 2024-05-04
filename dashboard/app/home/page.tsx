"use client"
import React, { useState } from 'react'
import { CSSProperties } from 'react';
import {
  QueryClient,
  QueryClientProvider,
  useQuery,
} from '@tanstack/react-query'
import CalendarStats from '@/components/CalendarStats'

const queryClient = new QueryClient()

const Page = () => {
  const [isHovered, setIsHovered] = useState(false);

  const onclick = async (e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    sessionStorage.removeItem('token');
    // Redirect to "/"
    window.location.href = "/";
  };

  const buttonStyle: CSSProperties = {
    position: 'fixed',
    top: '20px',
    right: '20px',
    width: 'auto',
    backgroundColor: '#ff6600',
    color: 'white',
    padding: '18px 24px',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
    fontSize: '18px',
    boxShadow: '0 4px 8px rgba(0, 0, 0, 0.3)',
    transition: 'transform 0.2s ease-in-out',
    transform: isHovered ? 'scale(1.05)' : 'none',
  };
  

  return (
    <div>
      <h1>Welcome to Home</h1>
      <button
        type="button"
        style={buttonStyle}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
        onClick={onclick} // Connect onclick function to onClick event
      >
        Logout
      </button>
    </div>
  );
};

export default Page;
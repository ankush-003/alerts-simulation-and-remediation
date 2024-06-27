"use client"
import React, { useState } from 'react'
import { CSSProperties } from 'react';
import {
  QueryClient,
  QueryClientProvider,
  useQuery,
} from '@tanstack/react-query'
import CalendarStats from '@/components/charts/CalendarStats'
import SeverityStats from '@/components/charts/SeverityStats';
import { toast } from 'sonner';

const queryClient = new QueryClient()

const Page = () => {
  const [isHovered, setIsHovered] = useState(false);

  const onclick = async (e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    sessionStorage.removeItem('token');
    // Redirect to "/"
    window.location.href = "/";
  };

  // toast.loading("Loading...", { duration: 2000 })

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
    <QueryClientProvider client={queryClient}>
      <div>
        {/* <h1>Welcome to Home</h1> */}
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
      <div className='w-full gap-2 grid'>
        <SeverityStats />
        <CalendarStats />
      </div>
    </QueryClientProvider>
  );
};

export default Page;
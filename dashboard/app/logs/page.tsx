"use client";
import React, { useState, useEffect } from 'react';

export default function Logs() {
  const [logs, setLogs] = useState([]);

  useEffect(() => {
    const fetchLogs = async () => {
      try {
        const jwtToken = sessionStorage.getItem('token');
        const response = await fetch('http://localhost:8000/alerts', {
          headers: {
            Authorization: `Bearer ${jwtToken}`,
          },
        });

        if (!response.ok) {
          throw new Error('Failed to fetch logs');
        }

        const data = await response.json();
        setLogs(data);
        console.log(data);
      } catch (error) {
        console.error('Error fetching logs:', error);
      }
    };

    fetchLogs();
  }, []);

  return (
    <div className='container mx-auto mt-8'>
      <h1 className='text-2xl font-bold mb-4'>Logs</h1>
      <div>
        {logs === null? (
          <p>No alert preferences have been selected.</p>
        ) : (
          logs.map((log, index) => (
            <div key={index} className="border border-gray-300 p-4 mb-4">
              <p>Acknowledged: {log.Acknowledged ? 'true' : 'false'}</p>
              <p>Category: {log.Category}</p>
              <p>CreatedAt: {log.CreatedAt}</p>
              <p>Remedy: {log.Remedy}</p>
              <p>Severity: {log.Severity}</p>
              <p>Source: {log.Source}</p>
              <p>Node: {log.node}</p>
              <p>_id: {log._id}</p>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
"use client";
import React, { useState, useEffect } from 'react';

interface Log {
  acknowledged: boolean;
  category: string;
  createdAt: string;
  remedy: string;
  severity: string;
  source: string;
  node: string;
  _id: string;
}

export default function Logs() {
  const [logs, setLogs] = useState<Log[] | null>(null);

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
              <p>Acknowledged: {log.acknowledged ? 'true' : 'false'}</p>
              <p>Category: {log.category}</p>
              <p>CreatedAt: {log.createdAt}</p>
              <p>Remedy: {log.remedy}</p>
              <p>Severity: {log.severity}</p>
              <p>Source: {log.source}</p>
              <p>Node: {log.node}</p>
              <p>_id: {log._id}</p>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
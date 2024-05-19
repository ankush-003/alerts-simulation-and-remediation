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

function formatTimestamp(timestamp: string): string {
  const date = new Date(timestamp);

  // Options for the date and time format
  const options: Intl.DateTimeFormatOptions = {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: 'numeric',
    minute: 'numeric',
    second: 'numeric',
  };

  // Use toLocaleDateString for a readable format
  return date.toLocaleDateString('en-US', options);
}

export default function Logs() {
  const [logs, setLogs] = useState<Log[] | null>(null);
  const [filteredLogs, setFilteredLogs] = useState<Log[] | null>(logs);
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(10);
  const [showAcknowledged, setShowAcknowledged] = useState(false);

  const fetchLogs = async (page: number, limit: number) => {
    try {
      const jwtToken = sessionStorage.getItem('token');
      const response = await fetch(`http://localhost:8000/alerts?page=${page}&limit=${limit}`, {
        headers: {
          Authorization: `Bearer ${jwtToken}`,
        },
      });

      if (!response.ok) {
        throw new Error('Failed to fetch logs');
      }

      const data = await response.json();
      setLogs(data.alerts);
      console.log(data);
    } catch (error) {
      console.error('Error fetching logs:', error);
    }
  };

  useEffect(() => {
    fetchLogs(page, limit);
  }, [page, limit]);

  useEffect(() => {
    setFilteredLogs(showAcknowledged
      ? logs?.filter(log => log.acknowledged === false) || []
      : logs || [])
  }, [showAcknowledged, logs])

  const handleNextPage = () => {
    setPage(prevPage => prevPage + 1);
  };

  const handlePrevPage = () => {
    setPage(prevPage => (prevPage > 1 ? prevPage - 1 : 1));
  };

  const toggleAcknowledged = () => {
    setShowAcknowledged(prevState => !prevState);
  };


  const acknowledgeLog = async (id: string) => {
    try {
      const jwtToken = sessionStorage.getItem('token');
      const response = await fetch(`http://localhost:8000/acknowledge?id=${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${jwtToken}`,
        },
        body: JSON.stringify({ id }),
      });

      if (!response.ok) {
        throw new Error('Failed to acknowledge log');
      }

      // Update the state to reflect the acknowledged log
      setLogs(prevLogs =>
        prevLogs?.map(log =>
          log._id === id ? { ...log, acknowledged: true } : log
        ) || null
      );
    } catch (error) {
      console.error('Error acknowledging log:', error);
    }
  };
  

  return (
    <div className='container mx-auto mt-8'>
      <h1 className='text-2xl font-bold mb-4'>Logs</h1>
      <button 
        onClick={toggleAcknowledged} 
        className="bg-blue-500 text-white py-2 px-4 rounded mb-4"
      >
        {showAcknowledged ? 'Show All Logs' : 'Show not Acknowledged only'}
      </button>
      <>
        {filteredLogs === null ? (
          <p>No alert preferences have been selected or no logs to display.</p>
        ) : (
          filteredLogs.map((log, index) => (
            <div key={index} className="border border-gray-300 p-4 mb-4">
              <p>Acknowledged: {log.acknowledged ? 'true' : 'false'}</p>
              <p>Category: {log.category}</p>
              <p>CreatedAt: {formatTimestamp(log.createdAt)}</p>
              <p>Remedy: {log.remedy}</p>
              <p>Severity: {log.severity}</p>
              <p>Source: {log.source}</p>
              <p>Node: {log.node}</p>
              <button 
                onClick={() => acknowledgeLog(log._id)} 
                disabled={log.acknowledged}
                className="bg-blue-500 text-white py-2 px-4 rounded mb-4"
              >
                Acknowledge
              </button>
            </div>
          ))
        )}
      </>
      <div className="flex justify-between mt-4">
        <button 
          onClick={handlePrevPage} 
          disabled={page === 1} 
          className="bg-blue-500 text-white py-2 px-4 rounded"
        >
          Previous
        </button>
        <button 
          onClick={handleNextPage} 
          className="bg-blue-500 text-white py-2 px-4 rounded"
        >
          Next
        </button>
      </div>
    </div>
  );
}

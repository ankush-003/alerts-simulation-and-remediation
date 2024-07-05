"use client";
import React, { useState, useEffect } from "react";
import { useSession } from "next-auth/react";
import { toast } from "sonner";
import { Button } from "./ui/button";
import { cn } from "@/lib/utils";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"

const backendUrl =
  process.env.NEXT_PUBLIC_BACKEND_URL || "http://localhost:8000";

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
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "numeric",
    minute: "numeric",
    second: "numeric",
  };

  // Use toLocaleDateString for a readable format
  return date.toLocaleDateString("en-US", options);
}

export default function Logs({ session: { user } }: any) {
  const [logs, setLogs] = useState<Log[] | null>(null);
  const [filteredLogs, setFilteredLogs] = useState<Log[] | null>(logs);
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(10);
  const [showAcknowledged, setShowAcknowledged] = useState(false);
  //   const { data: session } = useSession();

  const fetchLogs = async (page: number, limit: number) => {
    const toastId = toast.loading("Fetching logs...");
    try {
      const jwtToken = sessionStorage.getItem("token");
      const response = await fetch(
        `${backendUrl}/alerts?page=${page}&limit=${limit}`,
        {
          headers: {
            Authorization: `Bearer ${user?.token}`,
          },
        }
      );

      if (!response.ok) {
        throw new Error("Failed to fetch logs");
      }

      const data = await response.json();
      setLogs(data.alerts);
      toast.success("Logs fetched successfully", { id: toastId });
      //   console.log(data);
    } catch (error) {
      console.error("Error fetching logs:", error);
      toast.error("Failed to fetch logs", { id: toastId });
    }
  };

  useEffect(() => {
    fetchLogs(page, limit);
  }, [page, limit]);

  useEffect(() => {
    setFilteredLogs(
      showAcknowledged
        ? logs?.filter((log) => log.acknowledged === false) || []
        : logs || []
    );
  }, [showAcknowledged, logs]);

  const handleNextPage = () => {
    setPage((prevPage) => prevPage + 1);
  };

  const handlePrevPage = () => {
    setPage((prevPage) => (prevPage > 1 ? prevPage - 1 : 1));
  };

  const toggleAcknowledged = () => {
    setShowAcknowledged((prevState) => !prevState);
  };

  const acknowledgeLog = async (id: string) => {
    const toastId = toast.loading("Acknowledging log with id: " + id);
    try {
      const jwtToken = sessionStorage.getItem("token");
      const response = await fetch(`${backendUrl}/acknowledge?id=${id}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${user?.token}`,
        },
        body: JSON.stringify({ id }),
      });

      if (!response.ok) {
        toast.error("Failed to acknowledge log with id: " + id, {
          id: toastId,
        });
        throw new Error("Failed to acknowledge log");
      } else {
        toast.success(`alert with id: ${id} acknowledged successfully`, {
          id: toastId,
        });
      }

      // Update the state to reflect the acknowledged log
      setLogs(
        (prevLogs) =>
          prevLogs?.map((log) =>
            log._id === id ? { ...log, acknowledged: true } : log
          ) || null
      );
    } catch (error) {
      console.error("Error acknowledging log:", error);
    }
  };

  return (
    <div className="container mx-auto mt-8">
      <Button
        onClick={toggleAcknowledged}
        className="rounded mb-4"
        variant="secondary"
      >
        {showAcknowledged ? "Show All Logs" : "Show not Acknowledged only"}
      </Button>
      <>
        {filteredLogs === null ? (
            <div className="flex justify-center items-center">
                <p className="ml-4">No logs to display</p>
            </div>
        ) : (
          filteredLogs.map((log, index) => (
            <Alert key={index} className={
                cn(
                    log.acknowledged && "bg-green-500 text-white",
                    "p-4 mb-4"
                )
            }>
                <AlertTitle>{log.category}</AlertTitle>
                <AlertDescription>{log.remedy}</AlertDescription>
                <div className="flex justify-between">
                    <p>Source: {log.source}</p>
                    <p>Node: {log.node}</p>
                </div>
                <div className="flex justify-between">
                    <p>CreatedAt: {formatTimestamp(log.createdAt)}</p>
                    <Button
                    onClick={() => acknowledgeLog(log._id)}
                    disabled={log.acknowledged}
                    className={cn("mt-4", log.acknowledged && "bg-green-500")}
                    >
                    {log.acknowledged ? "Acknowledged" : "Acknowledge"}
                    </Button>
                </div>
            </Alert>
          ))
        )}
      </>
      <div className="flex justify-between mt-4">
        <Button
          onClick={handlePrevPage}
          disabled={page === 1}
          variant="ghost"
        >
          Previous
        </Button>
        <Button
          variant="ghost"
          onClick={handleNextPage}
        >
          Next
        </Button>
      </div>
    </div>
  );
}

"use client";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { Progress } from "@/components/ui/progress";
import { CardContent, Card } from "@/components/ui/card";
import { Computer, Cpu, MemoryStick } from "lucide-react";

interface HeartBeat {
  nodeID: string;
  numGoroutine: string;
  cpuUsage: string;
  ramUsage: string;
  status: string;
}

export default function Nodes() {
  const [nodes, setNodes] = useState<{ [key: string]: HeartBeat }>({});
  const [lastUpdatedTimes, setLastUpdatedTimes] = useState<{
    [key: string]: number;
  }>({});

  useEffect(() => {
    const eventSource = new EventSource("/api/stream/heartbeats");

    eventSource.addEventListener("message", (event) => {
      const data = JSON.parse(event.data);

      console.table(data);
      if (data.hasOwnProperty("message")) {
        // toast.success("Connection established");
        return;
      }
      setNodes((prevNodes) => ({
        ...prevNodes,
        [data.nodeID]: data,
      }));
      setLastUpdatedTimes((prevTimes) => ({
        ...prevTimes,
        [data.nodeID]: Date.now(),
      }));

      toast.success(`Heartbeat from ${data.nodeID}`);
    });

    eventSource.addEventListener("error", (event) => {
      console.error("SSE error:", event);
    });

    const checkDownStatus = setInterval(() => {
      setNodes((prevNodes) => {
        const updatedNodes = { ...prevNodes };
        Object.keys(prevNodes).forEach((nodeId) => {
          const lastUpdated = lastUpdatedTimes[nodeId];
          const fiveMinutesAgo = Date.now() - 1 * 60 * 1000;
          if (lastUpdated && lastUpdated < fiveMinutesAgo) {
            toast.error(`Node ${nodeId} is down`);
            updatedNodes[nodeId] = { ...prevNodes[nodeId], status: "down" };
          }
        });
        return updatedNodes;
      });
    }, 60 * 1000); // Check every minute

    return () => {
      eventSource.close();
      clearInterval(checkDownStatus);
    };
  }, [lastUpdatedTimes]);
  return (
    <>
      {Object.keys(nodes).length === 0 && (
        <div className="text-center">
          <h2 className="text-xl font-semibold text-center text-red-500">
            No nodes found
          </h2>
          <p className="text-lg text-center">
            Please make sure the nodes are running
          </p>
        </div>
      )}
      <section className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 p-4 md:p-6">
        {Object.keys(nodes).map((nodeId) => {
          const node = nodes[nodeId];
          const lastUpdated = lastUpdatedTimes[nodeId];
          console.table(node);
          const isDown = node.status === "down";
          return (
            <Card className="border-2 border-gray-200 p-4" key={node.nodeID}>
              <CardContent className="grid gap-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    {/* <Computer className="w-6 h-6 text-gray-500 dark:text-gray-400" /> */}
                    <h3 className="font-bold text-lg">
                      <span className={isDown ? "text-red-500" : "text-green-500"}>{node.nodeID}</span>
                    </h3>
                  </div>
                  {isDown ? (
                    <Computer className="w-6 h-6 text-red-500" />
                  ) : (
                    <Computer className="w-6 h-6 text-green-500" />
                  )}
                </div>
                <div className="grid gap-2">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <Cpu className="w-4 h-4 text-gray-500 dark:text-gray-400" />
                      <span className="text-gray-500 dark:text-gray-400">
                        CPU Usage
                      </span>
                    </div>
                    {node.cpuUsage === "-1" ?? (
                      <span className="font-medium text-red-500">
                        Not connected to Prometheus
                      </span>
                    )}
                  </div>
                  {node.cpuUsage === "-1" ? (
                    <Progress value={0} />
                  ) : (
                    <Progress value={parseFloat(node.cpuUsage)} />
                  )}
                </div>
                <div className="grid gap-2">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <MemoryStick className="w-4 h-4 text-gray-500 dark:text-gray-400" />
                      <span className="text-gray-500 dark:text-gray-400">
                        RAM Usage
                      </span>
                    </div>
                    {node.ramUsage === "-1" ?? (
                      <span className="font-medium text-red-500">
                        Not connected to Prometheus
                      </span>
                    )}
                  </div>
                  {node.ramUsage === "-1" ? (
                    <Progress value={0} />
                  ) : (
                    <Progress value={parseFloat(node.ramUsage)} />
                  )}
                </div>
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Cpu className="w-4 h-4 text-gray-500 dark:text-gray-400" />
                    <span className="text-gray-500 dark:text-gray-400">
                      Threads
                    </span>
                  </div>
                  {node.numGoroutine == "-1" ?? (
                    <span className="font-medium text-red-500">
                      Not connected to Prometheus
                    </span>
                  )}
                </div>
                {node.numGoroutine == "-1" ? (
                  <span className="font-medium">0</span>
                ) : (
                  <span className="font-medium">{node.numGoroutine}</span>
                )}
              </CardContent>
            </Card>
          );
        })}
      </section>
    </>
  );
}

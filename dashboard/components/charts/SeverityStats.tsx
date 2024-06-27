"use client";
import { useQuery } from "@tanstack/react-query";
import { Card, CardContent, CardHeader, CardTitle } from "../ui/card";

interface Data {
  severity: string;
  count: number;
}

export default function SeverityStats() {
  const { isPending, error, data } = useQuery<Data[]>({
    queryKey: ["severity_stats"],
    queryFn: async () => {
      const response = await fetch("/api/stats?type=severity");
      const data: any[] = await response.json();
      return data;
    },
  });

  if (isPending) {
    // return <div className="font-semibold text-center">Loading...</div>;
  }

  if (error) {
    return <div className="text-red-500 text-center font-semibold">Error: {error.message}</div>;
  }

  return (
  <div className="grid grid-cols-2 gap-4">
    {data?.map((item) => (
      <Card key={item.severity} className={
        `${item.severity === "Critical" ? "hover:text-red-500" :
            item.severity === "Severe" ? "hover:text-orange-500" :
            item.severity === "Warning" ? "hover:text-yellow-500" :
            item.severity === "Safe" ? "hover:text-green-500" : "hidden"}`
    
      }>
        <CardHeader>
          <CardTitle>{item.severity}</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-center text-3xl font-semibold">{item.count}</div>
        </CardContent>
      </Card>
    ))}
  </div>);
}

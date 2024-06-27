"use client";
import { useQuery } from "@tanstack/react-query";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  ResponsiveContainer,
  Tooltip,
  Label,
  CartesianGrid,
  AreaChart,
  Area,
  Line,
} from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "../ui/card";
import { toast } from "sonner";

interface Data {
  day: string;
  value: number;
}

export default function CalendarStats() {
  const { isPending, error, data } = useQuery<Data[]>({
    queryKey: ["calendar_stats"],
    queryFn: async () => {
      const response = await fetch("/api/stats?type=calendar");
      return response.json();
    },
  });

  if (isPending) {
    // toast.info("Loading...");
  }

  if (error) {
    return <div>Error: {error.message}</div>;
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Alerts Timeline</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="w-full h-">
          <ResponsiveContainer width="100%" height={400}>
            <AreaChart data={data} width={600} height={400}>
              <defs>
                <linearGradient id="valueColour" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#ef4444" stopOpacity={0.8} />
                  <stop offset="95%" stopColor="#ef4444" stopOpacity={0} />
                </linearGradient>
              </defs>
              {/* <CartesianGrid /> */}
              <XAxis dataKey="day" name="date" />
              <YAxis dataKey="value" name="count" />
              <Tooltip />
              <Area
                dataKey={"value"}
                stroke="#ef4444"
                name="Alerts"
                fillOpacity={1}
                fill="url(#valueColour)"
              />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      </CardContent>
    </Card>
  );
}

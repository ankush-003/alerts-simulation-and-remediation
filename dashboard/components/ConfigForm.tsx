"use client";
import React, { useState } from "react";
import { Switch } from "@/components/ui/switch";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import UpdateConfig from "@/actions/Config";
import { toast } from "sonner";
import { useSession } from "next-auth/react";

const allCategories = [
  "Memory",
  "CPU",
  "Disk",
  "Network",
  "Security",
  "RuntimeMetrics",
  "Power",
];

const allSeverities = ["Warning", "Severe", "Critical"];

export default function ConfigForm({
  userAlerts,
  token,
}: {
  userAlerts: any;
  token: string | undefined;
}) {
  const {data: session, update} = useSession();
  const [categories, setCategories] = useState(userAlerts?.Categories);
  const [severities, setSeverities] = useState(userAlerts?.Severities);

  return (
    <div>
      <form
        className="flex flex-col gap-3"
        action={async (formData) => {
          const alerts = {
            severities: allSeverities.filter((severity) =>
              formData.get(severity)
            ),
            categories: allCategories.filter((category) =>
              formData.get(category)
            ),
          };
          // console.log(alerts);
          const res = await UpdateConfig(alerts, token);
          const toastId = toast.loading("Updating alerts...");
          if (res.ok) {
            toast.success("Alerts updated successfully", { id: toastId });
            update({user: {Alert: alerts}})
          } else {
            toast.error(`Failed to update alerts: ${res.error}`, {
              id: toastId,
            });
          }
        }}
      >
        <Card>
          <CardHeader>
            <CardTitle>Severities</CardTitle>
          </CardHeader>
          <CardContent>
            <div>
              {allSeverities.map((severity) => (
                <div
                  key={severity}
                  className="flex items-center justify-between"
                >
                  <label htmlFor={severity}>{severity}</label>
                  <input
                    type="checkbox"
                    id={severity}
                    name={severity}
                    checked={severities.includes(severity)}
                    onChange={(e) => {
                      setSeverities((prev: any) => {
                        if (prev.includes(severity)) {
                          return prev.filter((s:any) => s !== severity);
                        } else {
                          return [...prev, severity];
                        }
                      });
                    }}
                  />
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Categories</CardTitle>
          </CardHeader>
          <CardContent>
            <div>
              {allCategories.map((category) => (
                <div
                  key={category}
                  className="flex items-center justify-between"
                >
                  <label htmlFor={category}>{category}</label>
                  <input
                    type="checkbox"
                    id={category}
                    name={category}
                    checked={categories.includes(category)}
                    onChange={(e) => {
                      setCategories((prev: any) => {
                        if (prev.includes(category)) {
                          return prev.filter((c: any) => c !== category);
                        } else {
                          return [...prev, category];
                        }
                      });
                    }}
                  />
                </div>
              ))}
            </div>
          </CardContent>
        <CardFooter className="flex justify-end mt-1">
          <Button type="submit">Save</Button>
        </CardFooter>
        </Card>
      </form>
    </div>
  );
}

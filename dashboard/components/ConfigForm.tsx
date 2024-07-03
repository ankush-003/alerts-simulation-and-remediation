"use client";
import React from 'react'
import { Switch } from "@/components/ui/switch";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import UpdateConfig from '@/actions/Config';
import { toast } from 'sonner';

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

export default function ConfigForm({ userAlerts, token }: { userAlerts: any, token: string | undefined }) {
  return (
    <div>
      <form className='flex flex-col gap-3' action={async (formData) => {
        const alerts = {
          severities: allSeverities.filter((severity) => formData.get(severity)),
          categories: allCategories.filter((category) => formData.get(category)),
        };
        console.log(alerts);
        const res = await UpdateConfig(alerts, token);
      }
      }>
        <Card>
          <CardHeader>
            <CardTitle>Severities</CardTitle>
          </CardHeader>
          <CardContent>
            <div>
            {allSeverities.map((severity) => (
              <div key={severity} className="flex items-center justify-between">
                <label htmlFor={severity}>{severity}</label>
                <input type="checkbox" id={severity} name={severity} checked={userAlerts?.severities.includes(severity)} />
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
              <div key={category} className="flex items-center justify-between">
                <label htmlFor={category}>{category}</label>
                <input type="checkbox" id={category} name={category} checked={userAlerts?.categories.includes(category)} />
              </div>
            ))}
            </div>
          </CardContent>
        </Card>
        <CardFooter className='flex justify-end mt-1'>
          <Button type="submit">Save</Button>
        </CardFooter>
      </form>
    </div>
  )
}

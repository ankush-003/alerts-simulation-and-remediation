import React from 'react'
import AuthProvider from '@/contexts/AuthProvider'
import Logs from '@/components/Logs'
import { auth } from '@/auth'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export default async function page() {
  const session = await auth();
  return (
    <>
      <AuthProvider>
        <Card>
          <CardHeader>
            <CardTitle>Logs</CardTitle>
          </CardHeader>
          <CardContent>
            <Logs session={session} />
          </CardContent>
        </Card>
      </AuthProvider>
    </>
  )
}

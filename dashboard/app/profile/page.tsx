import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useSession } from "next-auth/react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { auth } from "@/auth";
import ConfigForm from "@/components/ConfigForm";
import AuthProvider from "@/contexts/AuthProvider";

export default async function Profile() {
  const session = await auth();
  return (
    <div className="flex flex-col justify-center items-center gap-x-4">
      <Card className="w-full">
        <CardHeader>
          <CardTitle>{`Hello! ${session?.user?.first_name}`}</CardTitle>
          <CardDescription>Personal Details</CardDescription>
        </CardHeader>
        <CardContent>
          <form>
            <div className="grid w-full items-center gap-4 justify-center">
              <div className="flex gap-4 w-full">
                <div className="flex flex-col space-y-1.5 w-full">
                  <Label htmlFor="name">First Name</Label>
                  <Input disabled id="name" value={session?.user?.first_name} />
                </div>
                <div className="flex flex-col space-y-1.5">
                  <Label htmlFor="name">Last Name</Label>
                  <Input disabled id="name" value={session?.user?.last_name} />
                </div>
              </div>
              <div className="flex flex-col space-y-1.5">
                <Label htmlFor="email">Email</Label>
                <Input disabled id="email" value={session?.user?.email} />
              </div>
            </div>
          </form>
        </CardContent>
      </Card>
      <Card className="w-full mt-4">
        <CardHeader>
          <CardTitle>Your Alert Configuration</CardTitle>
          <CardDescription>Configure your alert settings</CardDescription>
        </CardHeader>
        <CardContent>
          <AuthProvider>
            <ConfigForm
              userAlerts={session?.user?.Alert}
              token={session?.user?.token}
            />
          </AuthProvider>
        </CardContent>
      </Card>
    </div>
  );
}

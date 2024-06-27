import { Button } from "./ui/button";
import { auth } from "@/auth";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "./ui/dropdown-menu";
import { SignIn, SignOut } from "@/components/auth-components";
import Link from "next/link";
import { CircleUser } from 'lucide-react';

export default async function UserButton() {
  const session = await auth();
  if (!session?.user)
    return (
      <Link href="/" className="underline">
        <Button>Sign In</Button>
      </Link>
    );
  return (
    <div className="flex gap-2 items-center">
      <span className="hidden text-sm sm:inline-flex">
        {`Welcome, ${session.user?.first_name}`}
      </span>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
            <CircleUser className="relative w-8 h-8 rounded-full" />
        </DropdownMenuTrigger>
        <DropdownMenuContent className="w-56" align="end" forceMount>
          <DropdownMenuLabel className="font-normal">
            <div className="flex flex-col space-y-1">
              <p className="text-sm font-medium leading-none">
                {session.user?.first_name} {session.user?.last_name}
              </p>
              <p className="text-xs leading-none text-muted-foreground">
                {session.user?.email}
              </p>
            </div>
          </DropdownMenuLabel>
          <DropdownMenuItem>
            <SignOut />
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
}
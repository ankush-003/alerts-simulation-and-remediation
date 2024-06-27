"use server";
import { signIn } from "@/auth";
import { CredentialsSignin } from "next-auth";
import { redirect } from 'next/navigation';

const credentialsLogin = async (email: string, password: string) => {
    let returnAddress = "/home";
  try {
    const res = await signIn("credentials", {
      email,
      password,
      redirectTo: "/home",
    });
  } catch (error: any) {
    const err = error as CredentialsSignin;
    returnAddress = "/";
  }

  redirect(returnAddress);
};


export default credentialsLogin;
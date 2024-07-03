"use server";
import { signIn } from "@/auth";
import { CredentialsSignin } from "next-auth";
import { redirect } from "next/navigation";

const credentialsLogin = async (email: string, password: string) => {
  try {
    const res = await signIn("credentials", {
      email,
      password,
      redirect: false,
    });
    console.log(res);

    return {
      ok: true,
      error: null,
    };
  } catch (error: any) {
    const err = error as CredentialsSignin;
    return {
      ok: false,
      error: "Invalid credentials. Please try again.",
    };
  }
};

export default credentialsLogin;
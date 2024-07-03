import NextAuth from "next-auth";
import Credentials from "next-auth/providers/credentials";
import { NextAuthConfig } from "next-auth";

const backendUrl =
  process.env.NEXT_PUBLIC_BACKEND_URL || "http://localhost:8000";

export const authConfig: NextAuthConfig = {
  providers: [
    Credentials({
      credentials: {
        email: { label: "Email", type: "text" },
        password: { label: "Password", type: "password" },
      },
      authorize: async (credentials, req) => {
        const email = credentials?.email as string;
        const password = credentials?.password as string;
        const res = await fetch(`${backendUrl}/users/login`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ email, password }),
        });
        const authData = await res.json();
        // console.log(authData);
        if (!res.ok) {
          return false;
        }
        const user = authData?.user;
        user.token = authData?.token;
        return user;
      },
    }),
  ],
  callbacks: {
    async jwt({ token, user }) {
      // console.log("jwt callback");
      if (user) {
        token.first_name = user.first_name;
        token.last_name = user.last_name;
        token.email = user.email;
        token.token = user.token;
      }
      return token;
    },
    async session({ session, token }) {
      // console.log("session callback");
      if (token) {
        session.user.first_name = token.first_name as string;
        session.user.last_name = token.last_name as string;
        session.user.email = token.email as string;
        session.user.token = token.token as string;
      }

      return session;
    },
  },

  pages: {
    // signIn: "/signin",
    // signOut: "/signout",
    // error: "/error",
    // newUser: "/new-user",
  },
  session: {
    strategy: "jwt",
  },
  theme: {
    colorScheme: "dark",
  },
};

export const { handlers, signIn, signOut, auth } = NextAuth(authConfig);
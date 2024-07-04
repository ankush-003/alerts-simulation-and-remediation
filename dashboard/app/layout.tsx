import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { cn } from "@/lib/utils";
import Sidebar from "@/components/Sidebar";
import { ThemeToggle } from "@/components/ThemeToggle";
import { Separator } from "@/components/ui/separator";
import { Toaster } from "@/components/ui/sonner";
import { SpeedInsights } from "@vercel/speed-insights/next";
import UserButton from "@/components/UserButton";
import ThemeProvider from "@/contexts/ThemeProvider";
import QueryProvider from "@/contexts/QueryProvider";
import AuthProvider from "@/contexts/AuthProvider";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "ASMR dashboard",
  description: "Dashboard app for alerts simulation",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={cn("min-h-screen w-full bg-white flex", inter.className, {
          "debug-screens": process.env.NODE_ENV === "development",
        })}
      >
        <AuthProvider>
          <ThemeProvider
            attribute="class"
            defaultTheme="dark"
            enableSystem
            disableTransitionOnChange
          >
            {/* Sidebar */}
            <Sidebar />
            {/* Main Page */}
            <div className="p-8 w-full bg-background text-foreground">
              <div className="flex justify-between items-center mb-4">
                <a href="/">
                  <h1
                    className="scroll-m-20 text-3xl font-extrabold tracking-tight xl:text-4xl text-transparent bg-gradient-to-r bg-clip-text from-red-500 via-orange-500 to-red-500
            animate-text"
                  >
                    Alerts Simulation & Remediation
                  </h1>
                </a>
                <div className="flex items-center space-x-2">
                  <UserButton />
                  <ThemeToggle />
                </div>
              </div>
              <Separator />
              <QueryProvider>
                <div className="mt-4 p-4">
                  <Toaster position="top-center" richColors />
                  {children}
                </div>
              </QueryProvider>
            </div>
          </ThemeProvider>
        </AuthProvider>
        <SpeedInsights />
      </body>
    </html>
  );
}

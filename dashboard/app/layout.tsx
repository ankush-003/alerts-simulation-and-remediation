import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { cn } from "@/lib/utils";
import { ThemeProvider } from "@/components/theme-provider";
import Sidebar from "@/components/Sidebar";
import { ThemeToggle } from "@/components/ThemeToggle";
import { Separator } from "@/components/ui/separator";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "dashboard",
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
        className={cn(
          "min-h-screen w-full bg-white flex",
          inter.className,
          {
            "debug-screens": process.env.NODE_ENV === "development",
          }
        )}
      >
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
                <h1 className="scroll-m-20 text-3xl font-extrabold tracking-tight xl:text-4xl text-transparent bg-gradient-to-r bg-clip-text from-red-500 via-orange-500 to-red-500
            animate-text">
                  Alerts Simulation & Remediation
                </h1>
              </a>
              <div className="flex items-center space-x-2">
                <ThemeToggle />
              </div>
            </div>
            <Separator />
            <div className="mt-4 p-4">
              {children}
            </div>
          </div>
        </ThemeProvider>
      </body>
    </html>
  );
}

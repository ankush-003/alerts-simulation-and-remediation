"use client";
import { TypewriterEffect } from "../components/ui/typewriter-effect";
import { AlertTriangle } from "lucide-react";
import { useEffect, useRef } from "react";
import Link from "next/link";
import { Button } from "@/components/ui/button"
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger, } from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import Lottie from "lottie-react";
import alertAnimation from "../public/alertAnimation.json";
import LoginRegisterForm from "@/components/Login";
import { Globe } from "@/components/ui/globe";

export default function Home() {
  const words = [
    { text: "Keep", },
    { text: "up", },
    { text: "with", },
    { text: "your", },
    { text: "Alerts.", className: "text-red-500 dark:text-red-500", },
  ];

  return (
    <div className="flex flex-col items-center justify-center mt-4">
      <p className="mb-4 text-4xl font-semibold text-center mt-5">
        Welcome to <span className="text-red-500 inline-block">Alerts Simulation & Remediation</span>
      </p>
      <div className="mt-2">
        <TypewriterEffect words={words} />
      </div>
      <div className="animation mb-6 mt-10">
        <Lottie animationData={alertAnimation} height={400} width={400} />
      </div>
      <div className="flex flex-col md:flex-row space-y-4 md:space-y-0 space-x-0 md:space-x-4 mt-10">
        <Dialog>
          <DialogTrigger asChild>
            <Button variant="secondary" className="w-40 h-15 rounded-xl text-white border border-black text-2xl">Login</Button>
          </DialogTrigger>
          <DialogContent className="flex flex-col items-center justify-center">
            <LoginRegisterForm />
          </DialogContent>
        </Dialog>
      </div>
      {/* <Globe /> */}
    </div>
  );
}
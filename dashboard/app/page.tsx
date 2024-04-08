"use client";
import { TypewriterEffect } from "../components/ui/typewriter-effect";
import { AlertTriangle } from "lucide-react";
import { useEffect, useRef } from "react";
import Link from "next/link";

import Lottie from "lottie-react";
import alertAnimation from "../public/alertAnimation.json";

export default function Home() {
  const words = [
    {
      text: "Keep",
    },
    {
      text: "up",
    },
    {
      text: "with",
    },
    {
      text: "your",
    },
    {
      text: "Alerts.",
      className: "text-red-500 dark:text-red-500",
    },
  ];

  return (
    <div className="flex flex-col items-center justify-center h-[40rem] ">
      <p className="mb-4 text-4xl font-semibold text-center">
        Welcome to <span className="text-red-500 inline-block">Alerts Simulation & Remediation</span>
      </p>
      <div className="mt-2">
        <TypewriterEffect words={words} />
      </div>
      <div className="animation mb-6 mt-10">
        <Lottie animationData={alertAnimation} height={400} width={400} />
      </div>
      <div className="flex flex-col md:flex-row space-y-4 md:space-y-0 space-x-0 md:space-x-4 mt-10">
        <button className="w-40 h-10 rounded-xl bg-black border dark:border-white border-transparent text-white text-lg">
          Join now
        </button>
        <Link href = "/signup">
        <button className="w-40 h-10 rounded-xl bg-white text-black border border-black  text-lg">
          Signup
        </button>
        </Link>

      </div>

    </div>
  );
}

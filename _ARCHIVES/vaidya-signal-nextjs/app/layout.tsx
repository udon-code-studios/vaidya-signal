import Footer from "@/components/Footer";
import "./globals.css";
import type { Metadata } from "next";
import { Inter } from "next/font/google";

export const metadata: Metadata = {
  title: "Vaidya Signal",
  description: "Created by Leo Battalora",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="bg-white text-gray-900">
        <div className="flex justify-center h-screen w-screen">
          <div className="flex flex-col justify-between items-center p-8 w-full min-h-screen max-w-screen-md">
            {children}
            <Footer />
          </div>
        </div>
      </body>
    </html>
  );
}

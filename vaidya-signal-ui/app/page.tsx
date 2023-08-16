import Image from "next/image";

export default function Home() {
  return (
    <div className="w-full flex flex-col gap-2">
      <p className="text-center text-5xl font-bold">Vaidya Signal</p>
      <div className="p-6 font-serif">
        <p className="pb-4">
          The Vaidya Signal is a specific bottom reversal divergence signal
          which is triggered when the following conditions are met:
        </p>
        <p className="pl-4">
          <strong>1.</strong> the current low* (L2) is lower than the previous low (L1)
        </p>
        <p className="pl-8 text-sm italic">
          *a low is defined as having 3 days before and after whith higher
          closes
        </p>
        <p className="pt-2 pl-4">
        <strong>2.</strong> MACD and RSI are both higher at L2 than they were at L1
        </p>
        <p className="pt-2 pl-4">
        <strong>3.</strong> volume at the L2 is lower than it was at L1
        </p>
      </div>
    </div>
  );
}

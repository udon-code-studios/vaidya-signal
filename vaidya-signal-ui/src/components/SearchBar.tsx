import { useState } from "react";

export default function Search({ autofocus = false }: { autofocus?: boolean }) {
  const [input, setInput] = useState<string>("");

  return (
    <div>
      <label className="relative block">
        <span className="absolute inset-y-0 left-0 flex items-center pl-4 opacity-75">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            aria-hidden="true"
            className="text-skin-base h-6 w-6"
          >
            <path
              fill="currentColor"
              d="M19.023 16.977a35.13 35.13 0 0 1-1.367-1.384c-.372-.378-.596-.653-.596-.653l-2.8-1.337A6.962 6.962 0 0 0 16 9c0-3.859-3.14-7-7-7S2 5.141 2 9s3.14 7 7 7c1.763 0 3.37-.66 4.603-1.739l1.337 2.8s.275.224.653.596c.387.363.896.854 1.384 1.367l1.358 1.392.604.646 2.121-2.121-.646-.604c-.379-.372-.885-.866-1.391-1.36zM9 14c-2.757 0-5-2.243-5-5s2.243-5 5-5 5 2.243 5 5-2.243 5-5 5z"
            >
            </path>
          </svg>
        </span>
        <input
          className="block w-full rounded-lg border border-skin-fill 
        border-opacity-40 bg-skin-fill py-3 pl-12
        pr-3 placeholder:italic placeholder:text-opacity-75 
        focus:border-skin-accent focus:outline-none"
          placeholder="Enter a ticker..."
          type="text"
          name="search"
          value={input}
          onChange={(e) => setInput(e.currentTarget.value)}
          autoComplete="off"
          autoFocus={autofocus}
          onKeyDown={(e) => {
            if (e.key === "Enter" && input !== "") {
              window.location.href = `/ticker/${input}`;
            }
          }}
        />
      </label>
    </div>
  );
}

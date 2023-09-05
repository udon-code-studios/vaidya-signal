import * as icons from "@/components/icons";

export default function Footer() {
  return (
    <div className="w-full flex flex-col gap-2">
      <div className="w-full flex justify-between items-end gap-8">
        <a href="https://udon.studio" target="_blank">
          <div>
            <icons.UdonLogoBgBlack class="w-14 sm:w-20" />
          </div>
        </a>

        <div className="flex gap-4 pr-4">
          <a
            href="https://github.com/udon-code-studios/vaidya-signal"
            target="_blank"
            className="hover:underline"
          >
            GitHub
          </a>
          <a
            href="mailto:leo.battalora@gmail.com"
            target="_blank"
            className="hover:underline"
          >
            Email
          </a>
          <a
            href="https://twitter.com/subpar_program"
            target="_blank"
            className="hover:underline"
          >
            Twitter
          </a>
        </div>
      </div>
      <div className="w-full border-2 border-current" />
      <p className="text-xs">
        Copyright © 2023 Udon Code Studios. All rights reserved.
      </p>
    </div>
  );
}

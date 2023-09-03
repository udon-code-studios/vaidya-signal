import { SITE } from "../config";
import Hr from "./Hr.tsx";
import LightDarkButton from "./LightDarkButton.tsx";
import SearchLink from "./SearchLink.tsx";

export default function Header(
  props: { active: "watchlist" | "email" | "search" | undefined },
) {
  return (
    <header>
      <div className="mx-auto flex max-w-3xl flex-col items-center justify-between sm:flex-row">
        <div className="flex w-full justify-between p-4 items-center py-8">
          <a
            href="/"
            className="py-1 text-xl font-semibold sm:static sm:text-2xl"
          >
            {SITE.title}
          </a>
          <div className="flex ml-0 mt-0 w-auto gap-x-6 gap-y-0">
            <a
              href="/watchlist"
              className={"hover:text-skin-accent self-end " +
                (props.active === "watchlist" ? "text-skin-accent" : "")}
            >
              Watchlist
            </a>
            <a
              href="/email"
              className={"hover:text-skin-accent " +
                (props.active === "email" ? "text-skin-accent" : "")}
            >
              Email
            </a>
            <SearchLink active={props.active === "search"} />
            <LightDarkButton />
          </div>
        </div>
      </div>
      <Hr noPadding={false} ariaHidden={true} />
    </header>
  );
}

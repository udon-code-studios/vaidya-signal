import Hr from "./Hr";
import Socials from "./Socials";

export default function Footer(props: { noMarginTop?: boolean }) {
  const currentYear = new Date().getFullYear();

  return (
    <footer className={`w-full ${props.noMarginTop ? "" : "mt-auto"}`}>
      <Hr noPadding={true} ariaHidden={true} />
      <div className="flex flex-col items-center justify-between py-6 sm:flex-row-reverse sm:py-4">
        <Socials centered />
        <div className="my-2 flex flex-col items-center whitespace-nowrap sm:flex-row">
          <span>Copyright &#169; {currentYear}</span>
          <span className="hidden sm:inline">&nbsp;|&nbsp;</span>
          <span>All rights reserved.</span>
        </div>
      </div>
    </footer>
  );
}

import { SOCIALS } from "../config";
import socialIcons from "../assets/socialIcons";

export default function Socials(props: { centered?: boolean }) {
  return (
    <div
      className={`flex-wrap justify-center gap-1 ${
        props.centered ? "flex" : ""
      }`}
    >
      {SOCIALS.filter((social) => social.active).map((social) => (
        <a
          href={social.href}
          target="_blank"
          className="hover:rotate-6 hover:text-skin-accent p-1"
          title={social.linkTitle}
        >
          <div dangerouslySetInnerHTML={{ __html: socialIcons[social.name] }} />
        </a>
      ))}
    </div>
  );

}

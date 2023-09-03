export default function Hr(props: { noPadding: boolean; ariaHidden: boolean }) {
  return (
    <div className={`max-w-3xl mx-auto ${props.noPadding ? "px-0" : "px-4"}`}>
      <hr className="border-skin-line" aria-hidden={props.ariaHidden} />
    </div>
  );
}

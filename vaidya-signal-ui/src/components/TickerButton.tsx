export default function Search({ ticker = "QQQ" }: { ticker: string }) {
  return (
    <a href="dog.com">
      <div className="w-full flex justify-center border p-2 rounded-lg py-4 sm:py-2 hover:text-skin-accent hover:font-bold">
        <p className="">{ticker}</p>
      </div>
    </a>
  );
}

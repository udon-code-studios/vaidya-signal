export function daysAgo(date: Date) {
  const today = new Date();
  const diff = today.getTime() - date.getTime();
  const days = Math.floor(diff / (1000 * 60 * 60 * 24));

  if (days === 0) return "today!";
  if (days === 1) return "yesterday";

  return `${days} days ago`;
}

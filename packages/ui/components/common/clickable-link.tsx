import Link from "next/link";

const ClickableLink = ({
  label,
  link,
  className = "text-primary hover:underline",
}: {
  label: string;
  link:
    | string
    | {
        pathname: string;
        query: { [key: string]: string | number };
      };
  className?: string;
}) => {
  if (!label || !link) {
    return <span>-</span>;
  }

  return (
    <Link href={link} className={className}>
      {label}
    </Link>
  );
};

export default ClickableLink; 
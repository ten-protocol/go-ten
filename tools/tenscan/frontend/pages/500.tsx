import { ErrorType } from "@/src/types/interfaces";
import Error from "./_error";

function Custom500Error({
  customPageTitle,
  message,
  showRedirectText,
  redirectText,
  err,
  redirectLink,
  children,
}: ErrorType) {
  return (
    <Error
      heading={"Oops! Something went wrong."}
      message={
        message ||
        "We're experiencing technical difficulties. Please try again later."
      }
      statusText={customPageTitle || `An Error occured`}
      statusCode={500}
      showRedirectText={showRedirectText || true}
      redirectText={redirectText || "Home Page"}
      err={err}
      redirectLink={redirectLink}
    >
      {children}
    </Error>
  );
}

export default Custom500Error;
